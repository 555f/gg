package metrics

import "github.com/555f/gg/pkg/gg"

var _ gg.Plugin = &Plugin{}

func init() {
	gg.RegisterPlugin(new(Plugin).Name(), func(ctx *gg.Context) gg.Plugin {
		return &Plugin{ctx: ctx}
	})
}
