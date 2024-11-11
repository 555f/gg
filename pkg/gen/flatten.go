package gen

import (
	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

type Paths []PathName

func (p Paths) String() string {
	var path string
	for i, v := range p {
		if i > 0 {
			path += "."
		}
		path += v.Name
	}
	return path
}

type FlattenPath struct {
	Var      *types.Var
	Path     *jen.Statement
	Paths    Paths
	Children []FlattenPath
	IsArray  bool
}

type PathName struct {
	Name     string
	JSONName string
}

func (n PathName) Value() string {
	return n.Name
}

func (n PathName) JSON() string {
	if n.JSONName == "" {
		return n.Name
	}
	return n.JSONName
}

type FlattenProcessor struct {
	allPaths    []FlattenPath
	currentPath []PathName
}

func (p *FlattenProcessor) varsByType(t any) (vars types.Vars) {
	switch t := t.(type) {
	case *types.Struct:
		for _, f := range t.Fields {
			vars = append(vars, f)
		}
		return vars
	case *types.Named:
		if t.System() {
			return nil
		}
		return p.varsByType(t.Type)
	}
	return nil
}

func (p *FlattenProcessor) flattenVar(v *types.Var) {
	var jsonName string
	if tag, err := v.SysTags.Get("json"); err == nil {
		jsonName = tag.Name
	}
	p.currentPath = append(p.currentPath, PathName{Name: v.Name, JSONName: jsonName})

	vars := p.varsByType(v.Type)

	hasChildren := len(vars) == 0

	if v.IsPointer && !hasChildren {
		currentPath := make([]PathName, len(p.currentPath))
		copy(currentPath, p.currentPath)

		path := jen.Id(currentPath[0].Value()).Do(func(s *jen.Statement) {
			for i := 1; i < len(currentPath); i++ {
				s.Dot(currentPath[i].Value())
			}
		})

		p.allPaths = append(p.allPaths, FlattenPath{
			Var:   v,
			Paths: currentPath,
			Path:  path,
		})
	}

	if hasChildren {
		var (
			children []FlattenPath
			isArray  bool
		)

		if s, ok := v.Type.(*types.Slice); ok {
			isArray = true
			for _, v := range p.varsByType(s.Value) {
				children = append(children, new(FlattenProcessor).Flatten(v)...)
			}
		}

		currentPath := make([]PathName, len(p.currentPath))
		copy(currentPath, p.currentPath)

		path := jen.Id(currentPath[0].Value()).Do(func(s *jen.Statement) {
			for i := 1; i < len(currentPath); i++ {
				s.Dot(currentPath[i].Value())
			}
		})

		p.allPaths = append(p.allPaths, FlattenPath{
			Var:      v,
			Paths:    currentPath,
			Children: children,
			Path:     path,
			IsArray:  isArray,
		})
	} else {
		for _, v := range vars {
			p.flattenVar(v)
		}
	}

	p.currentPath = p.currentPath[:len(p.currentPath)-1]
}

func (p *FlattenProcessor) Flatten(v *types.Var) []FlattenPath {
	p.flattenVar(v)
	return p.allPaths
}

func Flatten(v *types.Var) []FlattenPath {
	return (&FlattenProcessor{}).Flatten(v)
}
