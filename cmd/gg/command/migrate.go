package command

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/555f/gg/pkg/strcase"

	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

var (
	wdMigrateFile, pkgName, output              string
	openapiDoc, apiDoc, logging, client, server bool
)

type swipe struct {
	Named         *types.Named
	Type          *types.Interface
	MethodOptions map[string]map[string]any
	Namespace     string
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration from swipe3 to gg",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if pkgName == "" {
			pkgName = "migrate"
		}

		wd, err := filepath.Abs(wdMigrateFile)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cfg := &packages.Config{
			Dir:     wd,
			Context: context.TODO(),
			Mode:    packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
			Env:     os.Environ(),
		}
		var escaped []string

		for _, arg := range args {
			escaped = append(escaped, "pattern="+arg)
		}

		pkgs, err := packages.Load(cfg, escaped...)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		swipeOptions := map[string]*swipe{}

		for _, pkg := range pkgs {
			for _, syntax := range pkg.Syntax {
				ast.Inspect(syntax, func(node ast.Node) bool {
					switch t := node.(type) {
					case *ast.CallExpr:
						if sel, ok := t.Fun.(*ast.SelectorExpr); ok {
							if id, ok := sel.X.(*ast.Ident); ok && (id.Name == "gokit" || id.Name == "echo") {
								switch sel.Sel.Name {
								case "Interface":
									if len(t.Args) < 2 {
										return true
									}
									callExpr, ok := t.Args[0].(*ast.CallExpr)
									if !ok {
										return true
									}
									parenExpr, ok := callExpr.Fun.(*ast.ParenExpr)
									if !ok {
										return true
									}
									starExpr, ok := parenExpr.X.(*ast.StarExpr)
									if !ok {
										return true
									}
									iface := pkg.TypesInfo.TypeOf(starExpr.X)
									if iface != nil {
										opts, ok := swipeOptions[iface.String()]
										if !ok {
											ifaceNamed, _ := iface.(*types.Named)
											ifaceType, _ := iface.Underlying().(*types.Interface)
											opts = &swipe{
												Named:         ifaceNamed,
												Type:          ifaceType,
												MethodOptions: map[string]map[string]any{},
											}
											swipeOptions[iface.String()] = opts
										}
										if basic, ok := t.Args[1].(*ast.BasicLit); ok {
											v, _ := strconv.Unquote(basic.Value)
											opts.Namespace = v
										}
									}
								case "MethodOptions":
									mSel, ok := t.Args[0].(*ast.SelectorExpr)
									if !ok {
										return true
									}
									var ifaceExpr ast.Expr
									ifaceExpr, ok = mSel.X.(*ast.SelectorExpr)
									if !ok {
										ifaceExpr = mSel.X
									}
									iface := pkg.TypesInfo.TypeOf(ifaceExpr)
									if iface == nil {
										return true
									}
									opts, ok := swipeOptions[iface.String()]
									if !ok {
										ifaceNamed, _ := iface.(*types.Named)
										ifaceType, _ := iface.Underlying().(*types.Interface)
										opts = &swipe{
											Named:         ifaceNamed,
											Type:          ifaceType,
											MethodOptions: map[string]map[string]any{},
										}
										swipeOptions[iface.String()] = opts
									}
									for _, arg := range t.Args[1:] {
										if callExpr, ok := arg.(*ast.CallExpr); ok {
											if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
												if selIdent, ok := sel.X.(*ast.Ident); ok && (selIdent.Name == "gokit" || selIdent.Name == "echo") {
													if len(callExpr.Args) == 1 {
														if _, ok := opts.MethodOptions[mSel.Sel.Name]; !ok {
															opts.MethodOptions[mSel.Sel.Name] = map[string]any{}
														}
														var v any
														switch t := callExpr.Args[0].(type) {
														case *ast.CompositeLit:
															var values []string
															for _, elt := range t.Elts {
																if basic, ok := elt.(*ast.BasicLit); ok {
																	s, _ := strconv.Unquote(basic.Value)
																	values = append(values, s)
																}
															}
															v = values
														case *ast.BasicLit:
															switch t.Kind {
															case token.STRING:
																v, _ = strconv.Unquote(t.Value)
															}
														}
														opts.MethodOptions[mSel.Sel.Name][sel.Sel.Name] = v
													}
												}
											}
										}
									}
								}
							}
						}
					}
					return true
				})
			}
		}

		for _, opts := range swipeOptions {
			f := NewFile(pkgName)
			f.Comment(opts.Named.Obj().Name())
			f.Comment("@gg:\"http\"")
			f.Comment("@gg:\"middleware\"")
			if logging {
				f.Comment("@gg:\"logging\"")
			}
			if client {
				f.Comment("@http-client")
			}
			if server {
				f.Comment("@http-server")
			}
			if openapiDoc {
				f.Comment("@http-openapi")
			}
			if apiDoc {
				f.Comment("@http-apidoc")
			}

			f.Type().Id(opts.Named.Obj().Name()).InterfaceFunc(func(g *Group) {
				for i := 0; i < opts.Type.NumMethods(); i++ {
					method := opts.Type.Method(i)
					sig := method.Type().(*types.Signature)
					g.Comment(method.Name())

					var (
						httpMethod   = "GET"
						httpPath     = opts.Namespace
						headerVars   = map[string]string{}
						queryVars    = map[string]string{}
						wrapResponse string
					)

					if mopts, ok := opts.MethodOptions[method.Name()]; ok {
						if v, ok := mopts["RESTMethod"]; ok {
							httpMethod = v.(string)
						}
						if v, ok := mopts["RESTPath"]; ok {
							httpPath = httpPath + v.(string)
						}
						if v, ok := mopts["RESTHeaderVars"].([]string); ok {
							for i := 0; i < len(v); i += 2 {
								headerVars[v[i+1]] = v[i]
							}
						}
						if v, ok := mopts["RESTQueryVars"].([]string); ok {
							for i := 0; i < len(v); i += 2 {
								queryVars[v[i+1]] = v[i]
							}
						}
						if v, ok := mopts["RESTWrapResponse"].(string); ok {
							wrapResponse = v
						}
						//if v, ok := mopts["RESTQueryValues"].([]string); ok {
						//}
					}

					var isNoWrapResponse bool
					for i := 0; i < sig.Results().Len(); i++ {
						if named, ok := sig.Results().At(i).Type().(*types.Named); ok {
							if named.Obj().Name() == "error" {
								continue
							}
						}
						if sig.Results().At(i).Name() == "" {
							isNoWrapResponse = true
							break
						}
					}

					methodComment := fmt.Sprintf("@http-method:\"%s\"", httpMethod)
					if httpPath != "" {
						methodComment += fmt.Sprintf(" @http-path:\"%s\"", httpPath)
					}

					g.Commentf(methodComment)

					if wrapResponse != "" {
						g.Commentf("@http-wrap-response:\"%s\"", wrapResponse)
					} else if isNoWrapResponse {
						g.Commentf("@http-nowrap-response")
					}

					g.Id(method.Name()).Op("(").Line().CustomFunc(Options{}, func(g *Group) {
						for i := 0; i < sig.Params().Len(); i++ {
							param := sig.Params().At(i)
							if headerName, ok := headerVars[param.Name()]; ok {
								var isRequired bool
								if strings.HasPrefix(headerName, "!") {
									headerName = headerName[1:]
									isRequired = true
								}

								comment := "@http-type:\"header\""
								if headerName != "" {
									comment += fmt.Sprintf(" @http-name:\"%s\"", headerName)
								}

								if isRequired {
									comment += " @http-required"
								}

								g.Comment(comment).Line()
							}
							if queryName, ok := queryVars[param.Name()]; ok {
								comment := "@http-type:\"query\""
								if queryName != "" {
									comment += fmt.Sprintf(" @http-name:\"%s\"", queryName)
								}
								g.Comment(comment).Line()
							}
							g.Id(param.Name()).Add(qual(param.Type(), false)).Op(",").Line()
						}
					}).Op(")").
						Op("(").Line().
						CustomFunc(Options{}, func(g *Group) {
							for i := 0; i < sig.Results().Len(); i++ {
								param := sig.Results().At(i)
								name := param.Name()
								if name == "" {
									name = typeName(param.Type(), strconv.Itoa(i))
								}
								g.Id(name).Op(" ").Add(qual(param.Type(), false)).Op(",").Line()
							}
						}).
						Op(")")
				}
			})
			fpath, err := filepath.Abs(filepath.Join(wd, output, strings.ToLower(strcase.ToScreamingSnake(opts.Named.Obj().Name()))+".go"))
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			if err := f.Save(fpath); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}
	},
}

type QualFunc func(path, name string) *Statement

func typeName(t types.Type, postfix string) string {
	switch t := t.(type) {
	default:
		return ""
	case *types.Basic:
		return t.Name() + postfix
	case *types.Pointer:
		return typeName(t.Elem(), postfix)
	case *types.Array:
		return typeName(t.Elem(), postfix)
	case *types.Slice:
		return typeName(t.Elem(), postfix)
	case *types.Map:
		return typeName(t.Elem(), postfix)
	case *types.Named:
		if t.Obj().Name() == "error" {
			return "err"
		}
		return strcase.ToLowerCamel(t.Obj().Name())
	}
}

func pointerFn(isPtr bool) func(s *Statement) {
	return func(s *Statement) {
		if isPtr {
			s.Op("*")
		}
	}
}

func qual(t any, isPointer bool) *Statement {
	switch t := t.(type) {
	case *types.Pointer:
		return qual(t.Elem(), true)
	case *types.Interface:
		return Interface()
	case *types.Map:
		return Do(pointerFn(isPointer)).Map(qual(t.Key(), false)).Add(qual(t.Elem(), false))
	case *types.Array:
		return Do(pointerFn(isPointer)).Index(jen.Lit(t.Len())).Add(qual(t.Elem(), false))
	case *types.Slice:
		return Do(pointerFn(isPointer)).Index().Add(qual(t.Elem(), false))
	case *types.Var:
		return Do(pointerFn(isPointer)).Id(t.Name()).Add(qual(t.Type(), false))
	case *types.Tuple:
		var params []jen.Code
		for i := 0; i < t.Len(); i++ {
			v := t.At(i)

			st := jen.Id(v.Name())
			typ := v.Type()
			if s, ok := typ.(*types.Slice); ok {
				//if v.IsVariadic {
				//	st.Op("...")
				//} else {
				st.Index()
				//}
				typ = s.Elem()
			}
			st.Add(qual(typ, false))
			params = append(params, st)
		}
		return Params(params...)
	case *types.Signature:
		s := Add(qual(t.Params, false))
		if t.Results().Len() == 1 && t.Results().At(0).Name() == "" {
			s.Add(qual(t.Results().At(0), false))
		} else {
			s.Add(qual(t.Results(), false))
		}
		return s
	case *types.Basic:
		return Do(pointerFn(isPointer)).Id(t.Name())
	case *types.Named:
		s := Do(pointerFn(isPointer))
		if t.Obj().Pkg() == nil {
			s.Id(t.Obj().Name())
		} else {
			s.Qual(t.Obj().Pkg().Path(), t.Obj().Name())
		}
		return s
	case *types.Func:
		return Func().Id(t.Name()).Add(qual(t.Type(), false))
	}
	return nil
}

func init() {
	migrateCmd.PersistentFlags().StringVar(&wdMigrateFile, "work-dir", "", "working directory")
	migrateCmd.PersistentFlags().StringVar(&pkgName, "pkg", "", "package name")
	migrateCmd.PersistentFlags().StringVar(&output, "output", "", "output files path relative work-dir")
	migrateCmd.PersistentFlags().BoolVar(&openapiDoc, "openapi", false, "added openapi doc tag for interface")
	migrateCmd.PersistentFlags().BoolVar(&apiDoc, "apidoc", false, "added api doc tag for interface")
	migrateCmd.PersistentFlags().BoolVar(&logging, "logging", false, "added logging middleware for interface")
	migrateCmd.PersistentFlags().BoolVar(&client, "client", false, "added client for interface")
	migrateCmd.PersistentFlags().BoolVar(&server, "server", false, "added server for interface")

	rootCmd.AddCommand(migrateCmd)
}
