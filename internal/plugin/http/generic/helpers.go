package generic

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func makeRequestParam(parentParam, param *options.EndpointParam, importFn types.QualFunc) *Statement {
	fldName := param.FldNameUnExport
	if param.Parent != nil {
		fldName = param.Parent.FldNameUnExport + param.FldName
	}
	return Id(fldName).Add(types.Convert(param.Type, importFn))
}
