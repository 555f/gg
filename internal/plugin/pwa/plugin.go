package pwa

import (
	"context"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/555f/gg/internal/plugin/pwa/html2go"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/555f/gg/pkg/strcase"
	"golang.org/x/sync/errgroup"

	"github.com/dave/jennifer/jen"
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

	g, _ := errgroup.WithContext(context.Background())

	var mux sync.Mutex

	for _, name := range structNames {
		s := structMap[name]
		name := name

		if t, ok := s.Named.Tags.Get("pwa-view"); ok {
			g.Go(func() error {
				f, err := p.buildView(structMap, name, t.Value, s)
				if err != nil {
					return err
				}
				mux.Lock()
				files = append(files, f)
				mux.Unlock()
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return
}

func (p *Plugin) buildView(structMap map[string]*gg.Struct, name string, path string, s *gg.Struct) (*file.GoFile, error) {
	structPath := filepath.Join(p.ctx.Workdir, strings.Replace(s.Named.Pkg.Path, p.ctx.PkgPath, "", 1))
	templatePath, err := filepath.Abs(filepath.Join(structPath, path))
	if err != nil {
		return nil, errors.Error(err.Error(), token.Position{})
	}
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, errors.Error(err.Error(), token.Position{})
	}

	f := file.NewGoFile(p.ctx.Module, filepath.Join(structPath, strcase.ToSnake(s.Named.Name+"_render.go")))

	hg := html2go.NewHTML2Go(name, f.Qual, structMap)
	codes, _ := hg.Parse(string(data))

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
	return f, nil
}

func (p *Plugin) Dependencies() []string { return nil }
