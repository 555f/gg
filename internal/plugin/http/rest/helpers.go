package rest

import (
	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/gen"
	"github.com/555f/gg/pkg/types"

	. "github.com/dave/jennifer/jen"
)

func makeAddQueryParam(recvName string, parentParam, param *options.EndpointParam, qualFunc types.QualFunc, timeFormat string) *Statement {
	fldName := param.FldNameUnExport
	if parentParam != nil {
		fldName = parentParam.FldNameUnExport + param.FldName
	}
	paramID := Id(recvName).Dot("params").Dot(fldName)
	if param.Required {
		return Id("q").Dot("Add").Call(Lit(param.Name), gen.FormatValue(paramID, param.Type, qualFunc, timeFormat))
	} else {
		return If(Add(paramID).Op("!=").Nil()).Block(
			Id("q").Dot("Add").Call(Lit(param.Name), gen.FormatValue(Op("*").Add(paramID), param.Type, qualFunc, timeFormat)),
		)
	}
}

func makeRequestStructParam(parentParam, param *options.EndpointParam, importFn types.QualFunc) *Statement {
	fldName := param.FldNameUnExport
	if parentParam != nil {
		fldName = parentParam.FldNameUnExport + param.FldName
	}
	paramID := Id(fldName)
	if !param.Required {
		paramID.Op("*")
	}
	paramID.Add(types.Convert(param.Type, importFn))
	return paramID
}

func makeSetFunc(recvName, methodRequestName string, parentParam, param *options.EndpointParam, importFn types.QualFunc) *Statement {
	fldName := param.FldNameUnExport
	fnName := param.FldName
	if parentParam != nil {
		fldName = parentParam.FldNameUnExport + param.FldName
		fnName = parentParam.FldName + param.FldName
	}
	return Func().Params(Id(recvName).Op("*").Id(methodRequestName)).Id("Set" + fnName).
		Params(Id(fldName).Add(types.Convert(param.Type, importFn))).Op("*").Id(methodRequestName).BlockFunc(func(g *Group) {
		g.Id(recvName).Dot("params").Dot(fldName).Op("=").Op("&").Id(fldName)
		g.Return(Id(recvName))
	})
}
