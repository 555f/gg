package gg

type NodeSet map[string]struct{}
type PluginSet map[string]Plugin

type Dependency map[string]NodeSet

type Graph struct {
	plugins      PluginSet
	nodes        NodeSet
	dependencies Dependency
	dependents   Dependency
}

func (g *Graph) Leaves() []string {
	leaves := make([]string, 0)
	for node := range g.nodes {
		if _, ok := g.dependencies[node]; !ok {
			leaves = append(leaves, node)
		}
	}
	return leaves
}

func (g *Graph) HasDependency(child, parent string) bool {
	deps := g.Dependencies(child)
	_, ok := deps[parent]
	return ok
}

func (g *Graph) Dependencies(root string) NodeSet {
	if _, ok := g.nodes[root]; !ok {
		return nil
	}
	out := make(NodeSet)
	searchNext := []string{root}

	nextFn := func(node string) NodeSet {
		return g.dependencies[node]
	}

	for len(searchNext) > 0 {
		var discovered []string
		for _, node := range searchNext {
			for nextNode := range nextFn(node) {
				if _, ok := out[nextNode]; !ok {
					out[nextNode] = struct{}{}
					discovered = append(discovered, nextNode)
				}
			}
		}
		searchNext = discovered
	}
	return out
}

func (g *Graph) registerPlugin(plugin Plugin) error {
	g.nodes[plugin.Name()] = struct{}{}
	g.plugins[plugin.Name()] = plugin

	for _, child := range plugin.Dependencies() {
		if child == plugin.Name() {
			return ErrSelfReferential
		}
		if g.HasDependency(plugin.Name(), child) {
			return ErrCircularDependencies
		}

		g.nodes[child] = struct{}{}

		g.addNodeEdge(g.dependents, plugin.Name(), child)
		g.addNodeEdge(g.dependencies, child, plugin.Name())
	}
	return nil
}

func (g *Graph) sortedLayers() [][]string {
	var layers [][]string
	clonedGraph := g.clone()
	for {
		leaves := clonedGraph.Leaves()
		if len(leaves) == 0 {
			break
		}
		layers = append(layers, leaves)
		for _, leafNode := range leaves {
			clonedGraph.remove(leafNode)
		}
	}
	return layers
}

func (g *Graph) Sorted() (plugins []Plugin) {
	nodeCount := 0
	layers := g.sortedLayers()
	for _, layer := range layers {
		nodeCount += len(layer)
	}
	plugins = make([]Plugin, 0, nodeCount)
	for _, layer := range layers {
		for _, node := range layer {
			plugins = append(plugins, g.plugins[node])
		}
	}
	return plugins
}

func (g *Graph) remove(node string) {
	for dependent := range g.dependents[node] {
		removeFromDependency(g.dependencies, dependent, node)
	}
	delete(g.dependents, node)
	for dependency := range g.dependencies[node] {
		removeFromDependency(g.dependents, dependency, node)
	}
	delete(g.dependencies, node)
	delete(g.nodes, node)
}

func (g *Graph) addNodeEdge(d Dependency, key, node string) {
	nodes, ok := d[key]
	if !ok {
		nodes = make(NodeSet)
		d[key] = nodes
	}
	nodes[node] = struct{}{}
}

func (g *Graph) clone() *Graph {
	return &Graph{
		dependencies: copyDependency(g.dependencies),
		dependents:   copyDependency(g.dependents),
		nodes:        copyNodeset(g.nodes),
	}
}

func newGraph() *Graph {
	return &Graph{
		nodes:        make(NodeSet, 128),
		dependencies: make(Dependency, 128),
		dependents:   make(Dependency, 128),
		plugins:      make(map[string]Plugin, 128),
	}
}
