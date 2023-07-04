package config

import (
	"go/token"
	"path/filepath"

	"github.com/hashicorp/go-multierror"

	"github.com/555f/gg/internal/plugin/config/env"
	"github.com/555f/gg/internal/plugin/config/options"
	"github.com/555f/gg/pkg/errors"
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/gg"
)

type Plugin struct {
	ctx *gg.Context
}

func (p *Plugin) Name() string { return "config" }

func (p *Plugin) Exec() (files []file.File, errs error) {
	configOutput := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("output", "internal/config/config_loader.go"))
	docOutput := filepath.Join(p.ctx.Workdir, p.ctx.Options.GetStringWithDefault("doc-output", "docs/CONFIG.md"))

	f := file.NewGoFile(p.ctx.Module, configOutput)

	files = append(files, f)

	var mf *file.TxtFile

	for _, st := range p.ctx.Structs {
		opts, err := options.Decode(st)
		if err != nil {
			errs = multierror.Append(errs, errors.Error(err.Error(), token.Position{}))
		}
		env.GenConfig(opts)(f)
		if opts.MarkdownDoc {
			if mf == nil {
				mf = file.NewTxtFile(docOutput)
				files = append(files, mf)
			}
			env.GenMarkdownDoc(opts)(mf)
		}
		if opts.EnvsFile != "" {
			envsFileOutput := filepath.Join(p.ctx.Workdir, opts.EnvsFile)
			jbf := file.NewTxtFile(envsFileOutput)
			files = append(files, jbf)
			env.GenEnvsFile(opts)(jbf)
		}
	}
	return
}

func (p *Plugin) Dependencies() []string { return nil }
