package env

import (
	"github.com/555f/gg/internal/plugin/config/options"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/types"
	"github.com/555f/gg/pkg/typetransform"
	. "github.com/dave/jennifer/jen"
)

const multierrorPkg = "github.com/hashicorp/go-multierror"

func GenConfig(c options.Config) func(f *file.GoFile) {
	return func(f *file.GoFile) {
		stType := f.Import(c.PkgPath, c.Name)
		f.Func().Id(c.ConstructName).Params().Params(Id("c").Op("*").Do(stType), Id("errs").Error()).BlockFunc(func(g *Group) {
			g.Id("c").Op("=").Op("&").Do(stType).Values()

			walkFields(c.Fields, func(parent *options.ConfigField, field options.ConfigField) {
				var envName string
				code := Id("c")

				pathFields := resolvePathFields(parent)
				for i := len(pathFields) - 1; i >= 0; i-- {
					code = code.Dot(pathFields[i].FieldName)
					envName += pathFields[i].Name + "_"
				}

				envName += field.Name

				var isString bool
				if bt, ok := field.Type.(*types.Basic); ok {
					isString = bt.IsString()
				}

				lookupIf := g.If(List(Id("s"), Id("ok")).Op(":=").Qual("os", "LookupEnv").
					Call(Lit(envName)), Id("ok")).
					BlockFunc(func(g *Group) {
						if isString {
							g.Add(code).Dot(field.FieldName).Op("=").Id("s")
						} else {
							transCode, paramID, _ := typetransform.For(field.Type).
								SetAssignID(Id("v")).
								SetValueID(Id("s")).
								SetOp(":=").
								SetQualFunc(f.Import).
								SetErrStatements(
									Id("errs").Op("=").Qual(multierrorPkg, "Append").Call(Id("errs"), Qual("fmt", "Errorf").Call(Lit("env "+envName+" failed parse: %w"), Err())),
								).Parse()

							g.Add(transCode)

							g.Add(code).Dot(field.FieldName).Op("=").Add(paramID)
						}
						if !field.UseZero && field.Required {
							if field.Zero != "" {
								g.If(Add(code).Dot(field.FieldName).Op("==").Id(field.Zero)).Block(
									Id("errs").Op("=").Qual(multierrorPkg, "Append").Call(Id("errs"), Qual("errors", "New").Call(Lit("env "+envName+" empty"))),
								)
							}
						}
					})
				if field.Required {
					lookupIf.Else().Block(
						Id("errs").Op("=").Qual(multierrorPkg, "Append").Call(Id("errs"), Qual("errors", "New").Call(Lit("env "+envName+" not set"))),
					)
				}
			})

			g.Return()
		})
	}
}
