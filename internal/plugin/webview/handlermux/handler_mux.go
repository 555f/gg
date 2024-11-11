package handlermux

import "github.com/dave/jennifer/jen"

type HandlerMuxGenerator struct{}

func (*HandlerMuxGenerator) genDeclAt51() jen.Code {
	return jen.Null().Type().Id("HandleFunc").Func().
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("r").Op("*").Id("Request"),
		).
		Params(jen.Id("string"), jen.Id("error"))
}

func (*HandlerMuxGenerator) genDeclAt112() jen.Code {
	return jen.Null().Type().Id("handlerMux").Struct(jen.Id("mu").Qual("sync", "RWMutex"), jen.Id("m").Map(jen.Id("string")).Id("muxEntry"))
}

func (*HandlerMuxGenerator) genDeclAt180() jen.Code {
	return jen.Null().Type().Id("muxEntry").Struct(jen.Id("h").Id("HandleFunc"), jen.Id("pattern").Id("string"))
}

func (*HandlerMuxGenerator) genDeclAt242() jen.Code {
	return jen.Null().Var().Id("DefaultHandlerMux").Op("=").Op("&").Id("defaultHandlerMux")
}

func (*HandlerMuxGenerator) genDeclAt285() jen.Code {
	return jen.Null().Var().Id("defaultHandlerMux").Id("handlerMux")
}

func (*HandlerMuxGenerator) genBeforeFuncType() jen.Code {
	return jen.Null().Type().Id("BeforeFunc").Func().Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("r").Op("*").Id("Request")).Params(jen.Qual("context", "Context"), jen.Id("error"))
}

func (*HandlerMuxGenerator) genHandlerMuxOptionType() jen.Code {
	return jen.Null().Type().Id("handlerMuxOption").Struct(jen.Id("before").Index().Id("BeforeFunc"))
}

func (*HandlerMuxGenerator) genOptionType() jen.Code {
	return jen.Null().Type().Id("Option").Func().Params(jen.Op("*").Id("handlerMuxOption"))
}

func (*HandlerMuxGenerator) genFuncWithBefore() jen.Code {
	return jen.Func().Id("WithBefore").Params(jen.Id("before").Op("...").Id("BeforeFunc")).Params(jen.Id("Option")).Block(jen.Return().Func().Params(jen.Id("hmo").Op("*").Id("handlerMuxOption")).Block(jen.Id("hmo").Dot("before").Op("=").Id("append").Call(jen.Id("hmo").Dot("before"), jen.Id("before").Op("..."))))
}

func (*HandlerMuxGenerator) genRequestType() jen.Code {
	return jen.Null().Type().Id("Request").Struct(
		jen.Id("Meta").Map(jen.String()).Any(),
		jen.Id("Method").Id("string"),
		jen.Id("Params").Qual("encoding/json", "RawMessage"),
	)
}

func (*HandlerMuxGenerator) genFuncExecute() jen.Code {
	return jen.Func().Params(jen.Id("h").Op("*").Id("handlerMux")).Id("Execute").Params(jen.Id("payload").Index().Id("byte"), jen.Id("opts").Op("...").Id("Option")).Params(jen.Id("string"), jen.Id("error")).Block(jen.Id("h").Dot("mu").Dot("RLock").Call(), jen.Defer().Id("h").Dot("mu").Dot("RUnlock").Call(), jen.Id("opt").Op(":=").Op("&").Id("handlerMuxOption").Values(), jen.For(jen.List(jen.Id("_"), jen.Id("applyOpt")).Op(":=").Range().Id("opts")).Block(jen.Id("applyOpt").Call(jen.Id("opt"))), jen.Null().Var().Id("req").Op("*").Id("Request"), jen.If(jen.Id("err").Op(":=").Qual("encoding/json", "Unmarshal").Call(jen.Id("payload"), jen.Op("&").Id("req")), jen.Id("err").Op("!=").Id("nil")).Block(jen.Return().List(jen.Lit(""), jen.Id("err"))), jen.If(jen.List(jen.Id("e"), jen.Id("ok")).Op(":=").Id("h").Dot("m").Index(jen.Id("req").Dot("Method")), jen.Id("ok")).Block(jen.Id("parentCtx").Op(":=").Qual("context", "TODO").Call(), jen.For(jen.List(jen.Id("_"), jen.Id("before")).Op(":=").Range().Id("opt").Dot("before")).Block(jen.Null().Var().Id("err").Id("error"), jen.List(jen.Id("parentCtx"), jen.Id("err")).Op("=").Id("before").Call(jen.Id("parentCtx"), jen.Id("req")), jen.If(jen.Id("err").Op("!=").Id("nil")).Block(jen.Return().List(jen.Lit(""), jen.Id("err")))), jen.Return().Id("e").Dot("h").Call(jen.Id("parentCtx"), jen.Id("req"))), jen.Return().List(jen.Lit(""), jen.Id("nil")))
}

func (*HandlerMuxGenerator) genFuncRegister() jen.Code {
	return jen.Func().Params(
		jen.Id("h").Op("*").Id("handlerMux"),
	).Id("Register").Params(
		jen.Id("pattern").Id("string"),
		jen.Id("handle").Id("HandleFunc"),
	).Block(
		jen.Id("h").Dot("mu").Dot("Lock").Call(),
		jen.Defer().Id("h").Dot("mu").Dot("Unlock").Call(),
		jen.If(jen.Id("h").Dot("m").Op("==").Id("nil")).Block(
			jen.Id("h").Dot("m").Op("=").Id("make").Call(jen.Map(jen.Id("string")).Id("muxEntry")),
		),
		jen.Id("e").Op(":=").Id("muxEntry").Values(
			jen.Id("h").Op(":").Id("handle"),
			jen.Id("pattern").Op(":").Id("pattern"),
		),
		jen.Id("h").Dot("m").Index(jen.Id("pattern")).Op("=").Id("e"),
	)
}

func (*HandlerMuxGenerator) genFuncRegister2() jen.Code {
	return jen.Func().Id("Register").Params(jen.Id("pattern").Id("string"), jen.Id("handle").Id("HandleFunc")).Block(
		jen.Id("DefaultHandlerMux").Dot("Register").Call(jen.Id("pattern"), jen.Id("handle")),
	)
}

func (*HandlerMuxGenerator) genFuncExecute2() jen.Code {
	return jen.Func().Id("Execute").Params(jen.Id("payload").Index().Id("byte"), jen.Id("opts").Op("...").Id("Option")).Params(jen.Id("string"), jen.Id("error")).Block(jen.Return().Id("DefaultHandlerMux").Dot("Execute").Call(jen.Id("payload"), jen.Id("opts").Op("...")))
}

func (g *HandlerMuxGenerator) Generate() jen.Code {
	group := jen.NewFile("")
	group.Add(g.genBeforeFuncType())
	group.Add(g.genHandlerMuxOptionType())
	group.Add(g.genOptionType())
	group.Add(g.genFuncWithBefore())
	group.Add(g.genRequestType())
	group.Add(g.genDeclAt51())
	group.Add(g.genDeclAt112())
	group.Add(g.genDeclAt180())
	group.Add(g.genDeclAt242())
	group.Add(g.genDeclAt285())
	group.Add(g.genFuncExecute())
	group.Add(g.genFuncRegister())
	group.Add(g.genFuncExecute2())
	group.Add(g.genFuncRegister2())

	return group
}

func New() *HandlerMuxGenerator {
	return &HandlerMuxGenerator{}
}
