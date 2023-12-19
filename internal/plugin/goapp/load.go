package goapp

import (
	"strings"

	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
	"github.com/mantyr/rjson"
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
		case html.CommentNode:
			varPath := strings.TrimSpace(c.Data)
			if strings.HasPrefix(varPath, "$") {
				varPath = varPath[1:]
				if p := s.Type.Path(varPath); p != nil {

					t := p.Var.Type

					switch tt := t.(type) {
					case *types.Slice:
						t = tt.Value
					case *types.Array:
						t = tt.Value
					}

					if b, ok := t.(*types.Basic); ok {
						var fmtValCode jen.Code
						switch {
						case b.IsFloat():
							fmtValCode = jen.Lit("%f")
						case b.IsInteger():
							fmtValCode = jen.Lit("%d")
						case b.IsString():
							fmtValCode = jen.Lit("%s")
						}
						codes = append(codes, jen.Qual(appPkg, "Text").Call(
							jen.Qual("fmt", "Sprintf").Call(fmtValCode, jen.Id("c").Dot(varPath)),
						))
					}
				}
			}
		case html.ElementNode:
			var code *jen.Statement

			switch c.Data {
			default:
				code = jen.Qual(appPkg, strcase.ToCamel(c.Data)).Call()
			case "svg":
				code = jen.Id("SVG").Call()
			case "path":
				code = jen.Id("SVGPath").Call()
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
				switch a.Key {
				default:
					key := strcase.ToCamel(a.Key)
					if key == "Id" {
						key = "ID"
					}
					code.Dot(key).Call(jen.Lit(a.Val))
				case ":class":
					var classes map[string]interface{}

					rjson.Unmarshal([]byte(a.Val), &classes)

					code.Dot("Class").CallFunc(func(g *jen.Group) {
						var values []jen.Code
						for k, v := range classes {
							switch v := v.(type) {
							case string:
								values = append(values, jen.Lit(k).Op(":").Id(v))
							}
						}
						g.Id("classNames").Call(jen.Id("C").Values(values...))
					})
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

			if vRangeKey != "" {
				code = jen.Qual(appPkg, "Range").Call(jen.Id("c").Dot(vRangeSlice)).Dot("Slice").Call(
					jen.Func().Params(jen.Id(vRangeKey).Int()).Qual(appPkg, "UI").Block(
						jen.Return(code.Dot("Body").Call(load(f, structMap, s, c, fmtFunc)...)),
					),
				)
			} else {
				if codes := load(f, structMap, s, c, fmtFunc); len(codes) > 0 {
					code.Dot("Body").Call(codes...)
				}
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
			if strings.TrimSpace(c.Data) != "" {
				codes = append(codes, jen.Qual(appPkg, "Text").Call(jen.Lit(c.Data)))
			}
		}
	}
	return
}
