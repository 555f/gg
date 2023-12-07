package goapp

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
	"github.com/hashicorp/go-multierror"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "goapp" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	for _, s := range p.ctx.Structs {

		if t, ok := s.Named.Tags.Get("goapp-template"); ok {
			templatePath, err := filepath.Abs(filepath.Join(p.ctx.Workdir, strings.ReplaceAll(s.Named.Pkg.Path, p.ctx.Module.Path, ""), t.Value))
			if err != nil {
				errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				continue
			}
			fmt.Println(templatePath)
			_, err = os.ReadFile(templatePath)
			if err != nil {
				errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
				continue
			}
		}
	}
	return
}

func (p *Plugin) Dependencies() []string { return nil }
