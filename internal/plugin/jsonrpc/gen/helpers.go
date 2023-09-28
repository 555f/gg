package gen

import (
	"github.com/555f/gg/internal/plugin/jsonrpc/options"
	"github.com/555f/gg/pkg/strcase"
	"github.com/555f/gg/pkg/types"
)

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
