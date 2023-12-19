package gen

import (
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

func makeSetFunc(recvName, methodRequestName string, parentParam, param *options.EndpointParam, importFn types.QualFunc) *jen.Statement {
	fldName := param.FldNameUnExport
	fnName := param.FldName
	if parentParam != nil {
		fldName = parentParam.FldNameUnExport + param.FldName
		fnName = parentParam.FldName + param.FldName
	}
	return jen.Func().Params(jen.Id(recvName).Op("*").Id(methodRequestName)).Id("Set" + fnName).
		Params(jen.Id(fldName).Add(types.Convert(param.Type, importFn))).Op("*").Id(methodRequestName).BlockFunc(func(g *jen.Group) {
		g.Id(recvName).Dot("params").Dot(fldName).Op("=").Op("&").Id(fldName)
		g.Return(jen.Id(recvName))
	})
}

func makeRequestStructParam(parentParam, param *options.EndpointParam, importFn types.QualFunc) *jen.Statement {
	fldName := param.FldNameUnExport
	if parentParam != nil {
		fldName = parentParam.FldNameUnExport + param.FldName
	}
	paramID := jen.Id(fldName)
	if !param.Required {
		paramID.Op("*")
	}
	paramID.Add(types.Convert(param.Type, importFn))
	return paramID
}

func endpointFuncName(iface options.Iface, ep options.Endpoint) string {
	return strcase.ToLowerCamel(iface.Name+ep.MethodName) + "Endpoint"
}

func reqStructName(iface options.Iface, ep options.Endpoint) string {
	return strcase.ToLowerCamel(iface.Name+ep.MethodName) + "Req"
}

func respStructName(iface options.Iface, ep options.Endpoint) string {
	return strcase.ToLowerCamel(iface.Name+ep.MethodName) + "Resp"
}

func reqDecFuncName(iface options.Iface, ep options.Endpoint) string {
	return strcase.ToLowerCamel(iface.Name+ep.MethodName) + "ReqDec"
}

func respEncFuncName(iface options.Iface, ep options.Endpoint) string {
	return strcase.ToLowerCamel(iface.Name+ep.MethodName) + "RespEnc"
}

func clientStructName(iface options.Iface) string {
	return iface.Name + "Client"
}

func isNamedType(t any) (ok bool) {
	_, ok = t.(*types.Named)
	return
}

func normalizeVarName(name string) string {
	return strcase.ToLowerCamel(name)
}
