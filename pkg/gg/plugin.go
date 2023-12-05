package gg

import (
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/types"
)

var pluginFactories = map[string]PluginFactory{}

type PluginFactory func(ctx *Context) Plugin

type Plugin interface {
	Name() string
	Exec() (files []file.File, errs error)
	Dependencies() []string
}

type PluginAfterGen interface {
	OnAfterGen() error
}

type Struct struct {
	Named *types.Named
	Type  *types.Struct
}

type Interfaces []*Interface

type Interface struct {
	Named *types.Named
	Type  *types.Interface
}

type Tags struct {
	name  string
	items *types.Tags
}

func (pt *Tags) Tag(name string) (*types.Tag, bool) {
	prefix := pt.name
	if prefix != "" {
		prefix = prefix + "-"
	}
	return pt.items.Get(prefix + name)
}

type Options struct {
	m map[string]any
}

func (o Options) GetString(name string) string {
	v, _ := o.m[name].(string)
	return v
}

func (o Options) GetStringWithDefault(name, defaultValue string) string {
	if v, ok := o.m[name].(string); ok {
		return v
	}
	return defaultValue
}

func (o Options) GetBool(name string) bool {
	v, _ := o.m[name].(bool)
	return v
}

func (o Options) GetInt(name string) int {
	v, _ := o.m[name].(int)
	return v
}

type Context struct {
	pluginGraph      *Graph
	Version          string
	Workdir          string
	PkgPath          string
	Module           *types.Module
	Interfaces       []*Interface
	OthersInterfaces []*Interface
	Structs          []*Struct
	Options          Options
}

func (c *Context) Plugin(name string) Plugin {
	return c.pluginGraph.plugins[name]
}

func RegisterPlugin(name string, f PluginFactory) {
	if _, ok := pluginFactories[name]; ok {
		panic("plugin " + name + " has already registered")
	}
	pluginFactories[name] = f
}
