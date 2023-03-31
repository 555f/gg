package middleware

import "github.com/555f/gg/pkg/gg"

func init() {
	gg.RegisterPlugin(new(Plugin))
}
