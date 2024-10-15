package pwa

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
	"golang.org/x/net/html"
)

const appPkg = "github.com/maxence-charriere/go-app/v9/pkg/app"

const (
	condUnknownType = iota
	condIfType
	condElseType
	condElseIfType
)

type ifType struct {
	Type int
	Val  string
}

func load(f *file.GoFile, structMap map[string]*gg.Struct, s *gg.Struct, n *html.Node, fmtFunc func(jen.Code)) (codes []jen.Code) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.ElementNode:
			var code *jen.Statement

			switch c.Data {
			default:
				var (
					methods    []*types.Func
					isFoundTag bool
				)

				if ss, ok := structMap[c.Data]; ok {
					code = jen.Op("&").Do(f.Import(ss.Named.Pkg.Path, ss.Named.Name))
					code.Values()
					code = jen.Call(code)

					methods = ss.Named.Methods
					isFoundTag = true
				} else if p := s.Type.Path(strcase.ToLowerCamel(c.Data)); p != nil {

					switch tt := p.Type.(type) {
					case *types.Named:
						if iface := tt.Interface(); iface != nil {
							methods = iface.Methods
						} else {
							methods = tt.Methods
						}
					}
					code = jen.Id("c").Dot(strcase.ToLowerCamel(c.Data))
					isFoundTag = true
				}
				if isFoundTag {
					methodMap := make(map[string]*types.Func, len(methods))
					for _, m := range methods {
						methodMap[m.Name] = m
					}
					for _, a := range c.Attr {
						var isDynamic bool

						key := a.Key
						if strings.HasPrefix(key, ":") {
							isDynamic = true
							key = key[1:]
						}

						if isDynamic {
							code.Dot(strcase.ToCamel(key)).Call(jen.Id("c").Dot(a.Val))
						} else {
							fnName := strcase.ToCamel(key)
							if f, ok := methodMap[fnName]; ok {
								if f.Sig.Params.Len() > 0 {
									pt := f.Sig.Params[0]

									tt := pt.Type
									if t, ok := pt.Type.(*types.Slice); ok && pt.IsVariadic {
										tt = t.Value
									}

									switch t := tt.(type) {
									case *types.Basic:
										var val any
										switch {
										default:
											val = a.Val
										case t.IsBool():
											val, _ = strconv.ParseBool(a.Val)
										case t.IsFloat():
											val, _ = strconv.ParseFloat(a.Val, 64)
										case t.IsInteger():
											i, _ := strconv.ParseInt(a.Val, 10, 64)
											val = int(i)
										}
										code.Dot(fnName).Call(jen.Lit(val))
									}
								}
							}
						}
					}
					if codes := load(f, structMap, s, c, fmtFunc); len(codes) > 0 {
						code.Dot("Children").Call(codes...)
					}
					codes = append(codes, code)
					continue
				} else {
					code = jen.Qual(appPkg, strcase.ToCamel(c.Data)).Call()
				}
			case "svg":
				var buf bytes.Buffer
				html.Render(&buf, c)
				code = jen.Qual(appPkg, "Raw").Call(jen.Lit(buf.String()))
				codes = append(codes, code)
				continue
			}
			var (
				vRangeSlice string
				vRangeKey   string
				vIf         string
			)
			for _, a := range c.Attr {
				if strings.HasPrefix(a.Key, "aria-") || a.Key == "tabindex" {
					continue
				}
				if strings.HasPrefix(a.Key, ":") {
					attrName := a.Key[1:]
					if strings.HasPrefix(attrName, "on") {
						methodName := strcase.ToCamel(attrName[2:])
						code.Dot("On" + methodName).Call(jen.Id("c").Dot(a.Val))
					} else {
						code.Dot(strcase.ToCamel(attrName)).Call(jen.Id("c").Dot(a.Val))
					}
				} else {
					switch a.Key {
					default:
						key := strcase.ToCamel(a.Key)
						if key == "Id" {
							key = "ID"
						}
						code.Dot(key).Call(jen.Lit(a.Val))
					case "v-range":
						idx := strings.Index(a.Val, "in")
						if idx > -1 {
							vRangeKey = strings.TrimSpace(a.Val[:idx])
							vRangeSlice = strings.TrimSpace(a.Val[idx+2:])
						}
					case "v-if":
						vIf = a.Val
					}
				}
			}

			switch {
			default:
				if codes := load(f, structMap, s, c, fmtFunc); len(codes) > 0 {
					code.Dot("Body").Call(codes...)
				}
			case vRangeKey != "":
				code = jen.Qual(appPkg, "Range").Call(jen.Id("c").Dot(vRangeSlice)).Dot("Slice").Call(
					jen.Func().Params(jen.Id(vRangeKey).Int()).Qual(appPkg, "UI").Block(
						jen.Return(code.Dot("Body").Call(load(f, structMap, s, c, fmtFunc)...)),
					),
				)
			}

			if vIf != "" {
				var lastNc *html.Node
				for nc := c; nc != nil; nc = nc.NextSibling {
					if nc.Type != html.ElementNode {
						continue
					}
					var ifType ifType

					for _, a := range nc.Attr {
						ifType.Val = a.Val

						switch a.Key {
						case "v-if":
							ifType.Type = condIfType
						case "v-else":
							ifType.Type = condElseType
						case "v-else-if":
							ifType.Type = condElseIfType
						}
					}
					if ifType.Type == condUnknownType {
						lastNc = nc
						break
					}

					switch ifType.Type {
					case condIfType:
						code = jen.Qual(appPkg, "If").Call(jen.Id(ifType.Val), code)
					case condElseType:
						code = code.Dot("Else").Call(
							jen.Qual(appPkg, strcase.ToCamel(nc.Data)).Call().Dot("Body").Call(load(f, structMap, s, nc, fmtFunc)...),
						)
					case condElseIfType:
						code = code.Dot("ElseIf").Call(
							jen.Id(ifType.Val),
							jen.Qual(appPkg, strcase.ToCamel(nc.Data)).Call().Dot("Body").Call(load(f, structMap, s, nc, fmtFunc)...),
						)
					}
				}
				c = lastNc.PrevSibling
			}

			codes = append(codes, code)
		case html.TextNode:
			// var varPath string
			// 	for _, a := range c.Attr {
			// 		if a.Key == "name" {
			// 			varPath = a.Val
			// 			break
			// 		}
			// 	}
			// 	if varPath == "" {
			// 		continue
			// 	}
			// 	p := s.Type.Path(varPath)
			// 	if p != nil {
			// 		t := p.Var.Type

			// 		switch tt := t.(type) {
			// 		case *types.Slice:
			// 			t = tt.Value
			// 		case *types.Array:
			// 			t = tt.Value
			// 		}

			// 		switch tt := t.(type) {
			// 		case *types.Basic:
			// 			var valCode jen.Code
			// 			switch {
			// 			case tt.IsNumeric():
			// 				valCode = jen.Qual("fmt", "Sprint").Call(jen.Id("c").Dot(varPath))
			// 			case tt.IsString():
			// 				valCode = jen.Id("c").Dot(varPath)
			// 			}
			// 			codes = append(codes, jen.Qual(appPkg, "Text").Call(valCode))
			// 		case *types.Named:
			// 			code := jen.Id("c").Dot(varPath)
			// 			for _, a := range c.Attr {
			// 				if a.Key != "name" {
			// 					code.Dot(strcase.ToCamel(a.Key)).Call(jen.Lit(a.Val))
			// 				}
			// 			}
			// 			codes = append(codes, code)
			// 		}
			// 	}

			var v jen.Code

			data := strings.TrimSpace(c.Data)

			if data, ok := strings.CutPrefix(data, "{"); ok {
				i := strings.Index(data, "}")
				varPath := data[:i]
				if varPath != "" {
					p := s.Type.Path(varPath)
					if p != nil {
						t := p.Type
						switch tt := t.(type) {
						case *types.Slice:
							t = tt.Value
						case *types.Array:
							t = tt.Value
						}

						switch tt := t.(type) {
						case *types.Basic:
							var valCode jen.Code
							switch {
							case tt.IsNumeric():
								valCode = jen.Qual("fmt", "Sprint").Call(jen.Id("c").Dot(varPath))
							case tt.IsString():
								valCode = jen.Id("c").Dot(varPath)
							}
							codes = append(codes, jen.Qual(appPkg, "Text").Call(valCode))
							// case *types.Named:
							// 	code := jen.Id("c").Dot(varPath)
							// 	for _, a := range c.Attr {
							// 		if a.Key != "name" {
							// 			code.Dot(strcase.ToCamel(a.Key)).Call(jen.Lit(a.Val))
							// 		}
							// 	}
							// 	codes = append(codes, code)
						}
					}
				}
			} else if data != "" {
				v = jen.Lit(data)

				codes = append(codes, jen.Qual(appPkg, "Text").Call(v))
			}

		}
	}
	return
}
