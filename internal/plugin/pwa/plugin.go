package pwa

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/555f/gg/internal/plugin/pwa/html2go"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"

	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/go-multierror"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "pwa" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	structNames := make([]string, 0, len(p.ctx.Structs))
	structMap := make(map[string]*gg.Struct, len(p.ctx.Structs))

	for _, s := range p.ctx.Structs {
		name := strcase.ToKebab(s.Named.Name)
		structNames = append(structNames, name)
		structMap[name] = s
	}

	for _, name := range structNames {
		s := structMap[name]
		if t, ok := s.Named.Tags.Get("pwa-view"); ok {
			structPath := filepath.Join(p.ctx.Workdir, strings.Replace(s.Named.Pkg.Path, p.ctx.PkgPath, "", 1))
			templatePath, err := filepath.Abs(filepath.Join(structPath, t.Value))
			if err != nil {
				errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				continue
			}
			data, err := os.ReadFile(templatePath)
			if err != nil {
				errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				continue
			}

			f := file.NewGoFile(p.ctx.Module, filepath.Join(structPath, strcase.ToSnake(s.Named.Name+"_render.go")))

			hg := html2go.NewHTML2Go(name, f.Qual, structMap)
			codes, _ := hg.Parse(string(data))

			// reader := strings.NewReader(string(data))
			// doc, err := html.Parse(reader)
			// if err != nil {
			// errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
			// continue
			// }

			// codes := load(f, structMap, s, doc.FirstChild.FirstChild.NextSibling, func(c jen.Code) {
			// 	f.Add(c)
			// })

			f.Func().
				Params(jen.Id("c").Op("*").Id(s.Named.Name)).Id("Render").
				Params().Qual(appPkg, "UI").
				Block(
					jen.If(
						jen.List(jen.Id("cc"), jen.Id("ok")).
							Op(":=").
							Any().Call(jen.Id("c")).Assert(jen.Interface(jen.Id("OnBeforeRender").Call())),
						jen.Id("ok"),
					).Block(
						jen.Id("cc").Dot("OnBeforeRender").Call(),
					),
					jen.Return(
						codes...,
					),
				)

			files = append(files, f)
		}
	}
	return
}

func (p *Plugin) Dependencies() []string { return nil }
