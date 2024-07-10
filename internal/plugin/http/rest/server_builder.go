package rest

import (
	"sort"

	"github.com/555f/gg/internal/plugin/http/options"
	"github.com/555f/gg/pkg/strcase"
	"github.com/dave/jennifer/jen"
)

type BaseServerBuilder struct {
	typesStatement       jen.Code
	qualifier            Qualifier
	handlerStrategies    map[string]HandlerStrategyBuilderFactory
	errorWrapper         *options.ErrorWrapper
	loadedHandleStrategy map[string]HandlerStrategy
	codes                []jen.Code
}

func (b *BaseServerBuilder) Build() jen.Code {
	var statements []jen.Code
	var loadedHandleStrategyKeys []string

	for k := range b.loadedHandleStrategy {
		loadedHandleStrategyKeys = append(loadedHandleStrategyKeys, k)
	}
	sort.Strings(loadedHandleStrategyKeys)

	for _, k := range loadedHandleStrategyKeys {
		handleStrategy := b.loadedHandleStrategy[k]

		statements = append(statements,
			jen.Func().Id(handleStrategy.ID()+"DefaultErrorEncoder").Params(
				jen.Id(handleStrategy.RespArgName()).Add(handleStrategy.RespType()),
				jen.Err().Error(),
			).BlockFunc(func(g *jen.Group) {
				g.Var().Id("statusCode").Int().Op("=").Qual("net/http", "StatusInternalServerError")
				g.If(jen.List(jen.Id("e"), jen.Id("ok")).Op(":=").Err().Assert(jen.Interface(jen.Id("StatusCode").Params().Int())), jen.Id("ok")).Block(
					jen.Id("statusCode").Op("=").Id("e").Dot("StatusCode").Call(),
				)
				g.If(jen.List(jen.Id("headerer"), jen.Id("ok")).Op(":=").Err().Assert(jen.Interface(jen.Id("Headers").Params().Qual("net/http", "Header"))), jen.Id("ok")).Block(
					jen.For(jen.List(jen.Id("k"), jen.Id("values"))).Op(":=").Range().Id("headerer").Dot("Headers").Call().Block(
						jen.For(jen.List(jen.Id("_"), jen.Id("v"))).Op(":=").Range().Id("values").Block(

							jen.Add(handleStrategy.SetHeader(jen.Id("k"), jen.Id("v"))),
						),
					),
				)

				if b.errorWrapper != nil {
					errorWrapperName := strcase.ToLowerCamel(b.errorWrapper.Struct.Named.Name)
					g.Id(errorWrapperName).Op(":=").Do(b.qualifier.Qual(b.errorWrapper.Struct.Named.Pkg.Path, b.errorWrapper.Struct.Named.Name)).Values()
					for _, field := range b.errorWrapper.Fields {
						g.If(jen.List(jen.Id("e"), jen.Id("ok")).Op(":=").Err().Assert(jen.Interface(jen.Id(field.Interface))), jen.Id("ok")).Block(
							jen.Id(errorWrapperName).Dot(field.FldName).Op("=").Id("e").Op(".").Id(field.MethodName).Call(),
						)
					}
					g.Add(handleStrategy.WriteError(jen.Id("statusCode"), jen.Id(errorWrapperName)))
				} else {
					g.Add(handleStrategy.WriteError(jen.Id("statusCode"), jen.Id("err")))
				}
			}),

			jen.Func().Id("encodeBody").Params(
				jen.Id("rw").Do(b.qualifier.Qual(httpPkg, "ResponseWriter")),
				jen.Id("data").Any(),
			).Block(),
		)
	}

	if b.typesStatement != nil {
		statements = append(statements, b.typesStatement)
	}
	statements = append(statements, b.codes...)

	return jen.Custom(jen.Options{Multi: true}, statements...)
}

func (b *BaseServerBuilder) RegisterHandlerStrategy(name string, f HandlerStrategyBuilderFactory) {
	b.handlerStrategies[name] = f
}

func (b *BaseServerBuilder) Controller(iface options.Iface) ServerControllerBuilder {
	f, ok := b.handlerStrategies[iface.Type]
	if !ok {
		panic("unknown strategy " + iface.Type)
	}
	handlerStrategy := f()
	if _, ok := b.loadedHandleStrategy[iface.Type]; !ok {
		b.loadedHandleStrategy[iface.Type] = handlerStrategy
	}
	return &serverControllerBuilder{BaseServerBuilder: b, iface: iface, handlerStrategy: handlerStrategy}
}

func NewServerBuilder(qualifier Qualifier, errorWrapper *options.ErrorWrapper) *BaseServerBuilder {
	return &BaseServerBuilder{
		qualifier:            qualifier,
		handlerStrategies:    map[string]HandlerStrategyBuilderFactory{},
		loadedHandleStrategy: map[string]HandlerStrategy{},
		errorWrapper:         errorWrapper,
	}
}
