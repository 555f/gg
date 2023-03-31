package gg

import (
	"github.com/555f/gg/pkg/file"
	"github.com/555f/gg/pkg/types"
)

type Plugin interface {
	Name() string
	Exec(ctx *Context) (files []file.File, errs error)
	Dependencies() []string
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

func (pt *Tags) Get(name string) (*types.Tag, bool) {
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
	pluginGraph *Graph
	Module      *types.Module
	Interfaces  []*Interface
	Structs     []*Struct
	Options     Options
}

func (c *Context) Plugin(name string) Plugin {
	return c.pluginGraph.plugins[name]
}

type PluginFactory func() Plugin

var pluginGraph = newGraph()

func RegisterPlugin(plugin Plugin) {
	if err := pluginGraph.registerPlugin(plugin); err != nil {
		panic(err)
	}
}
