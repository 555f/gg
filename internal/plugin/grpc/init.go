package grpc

import "github.com/555f/gg/pkg/gg"

func init() {
	gg.RegisterPlugin(new(Plugin).Name(), func(ctx *gg.Context) gg.Plugin {
		return &Plugin{ctx: ctx}
	})
}
