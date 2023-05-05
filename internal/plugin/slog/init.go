package logging

import "github.com/555f/gg/pkg/gg"

var _ gg.Plugin = &Plugin{}

func init() {
	gg.RegisterPlugin(new(Plugin))
}
