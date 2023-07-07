package openapidoc

import (
	"github.com/555f/gg/internal/openapi"
	openapi2 "github.com/555f/gg/internal/openapi"
)

func wrapScheme(wrapPath []string, props openapi2.Properties) openapi.Properties {
	if len(wrapPath) == 0 {
		return props
	}
	name, wrapPath := wrapPath[0], wrapPath[1:]
	return openapi2.Properties{
		name: &openapi2.Schema{
			Type:       "object",
			Properties: wrapScheme(wrapPath, props),
		},
	}
}
