package types

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	stdtypes "go/types"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"

	"golang.org/x/tools/go/packages"
	stdpackages "golang.org/x/tools/go/packages"
)

type Decoder struct {
	pkg      *packages.Package
	packages []*stdpackages.Package
	visited  map[string]any
}

func (d *Decoder) Decode(t any) (any, error) {
	return d.normalizeRecursive(t, false)
}

func (d *Decoder) normalizeModule(module *packages.Module) (*Module, error) {
	if module != nil {
		return &Module{
			ID:       filepath.Base(module.Path),
			Version:  module.Version,
			Path:     module.Path,
			Dir:      module.Dir,
			Indirect: module.Indirect,
			Main:     module.Main,
		}, nil
	}
	return nil, nil
}

func (d *Decoder) normalizeBasic(t *stdtypes.Basic, isPointer bool) (*Basic, error) {

	return &Basic{
		Name:      t.Name(),
		IsPointer: isPointer,
		Kind:      t.Kind(),
		Zero:      zeroValue(t),
	}, nil
}

func (d *Decoder) normalizeChan(t *stdtypes.Chan) (any, error) {
	tp, err := d.normalizeRecursive(t.Elem(), false)
	if err != nil {
		return nil, err
	}
	return &Chan{Type: tp}, nil
}

func (d *Decoder) normalizeTuple(t *stdtypes.Tuple) (any, error) {
	tp := &Tuple{}
	for i := 0; i < t.Len(); i++ {
		//v, err := d.normalizeVar(t.At(i))
		//if err != nil {
		//	return nil, err
		//}
		//tp.Vars = append(tp.Vars, v)
	}
	return tp, nil
}

func (d *Decoder) normalizeStruct(t *stdtypes.Struct, isPointer bool) (*Struct, error) {
	result := &Struct{
		IsPointer: isPointer,
	}
	for i := 0; i < t.NumFields(); i++ {
		field := t.Field(i)
		v, err := d.normalizeVar(field)
		if err != nil {
			return nil, err
		}
		f := &StructFieldType{
			Var: v,
		}
		if tags, err := structtag.Parse(t.Tag(i)); err == nil {
			f.SysTags = tags
		}
		result.Fields = append(result.Fields, f)
	}
	return result, nil
}

func (d *Decoder) normalizePkg(pkg *stdtypes.Package) (*PackageType, error) {
	if pkg == nil {
		return nil, nil
	}
	module, err := d.normalizeModule(d.pkg.Module)
	if err != nil {
		return nil, err
	}
	return &PackageType{
		Name:   pkg.Name(),
		Path:   pkg.Path(),
		Module: module,
	}, nil
}

func (d *Decoder) normalizeMap(key stdtypes.Type, val stdtypes.Type, isPointer bool) (mapType *Map, err error) {
	mapType = &Map{IsPointer: isPointer}
	mapType.Key, err = d.normalizeRecursive(key, false)
	if err != nil {
		return nil, err
	}
	mapType.Value, err = d.normalizeRecursive(val, false)
	if err != nil {
		return nil, err
	}
	return mapType, nil
}

func (d *Decoder) normalizeSlice(val stdtypes.Type, isPointer bool) (*Slice, error) {
	v, err := d.normalizeRecursive(val, false)
	if err != nil {
		return nil, err
	}
	return &Slice{
		Value:     v,
		IsPointer: isPointer,
	}, nil
}

func (d *Decoder) normalizeArray(val stdtypes.Type, len int64, isPointer bool) (*Array, error) {
	v, err := d.normalizeRecursive(val, false)
	if err != nil {
		return nil, err
	}
	return &Array{
		Value:     v,
		Len:       len,
		IsPointer: isPointer,
	}, nil
}

func (d *Decoder) normalizeInterface(t *stdtypes.Interface) (*Interface, error) {
	it := &Interface{
		origin: t,
	}
	for i := 0; i < t.NumMethods(); i++ {
		method, err := d.normalizeFunc(t.Method(i))
		if err != nil {
			return nil, err
		}
		it.Methods = append(it.Methods, method)
	}
	for i := 0; i < t.NumEmbeddeds(); i++ {
		method, err := d.normalizeRecursive(t.EmbeddedType(i), false)
		if err != nil {
			return nil, err
		}
		it.Embedded = append(it.Embedded, method)
	}
	for i := 0; i < t.NumExplicitMethods(); i++ {
		method, err := d.normalizeFunc(t.ExplicitMethod(i))
		if err != nil {
			return nil, err
		}
		it.ExplicitMethods = append(it.ExplicitMethods, method)
	}
	return it, nil
}

func (d *Decoder) normalizeNamed(named *stdtypes.Named, isPointer bool) (nt *Named, err error) {
	var pkgPath string
	if named.Obj().Pkg() != nil {
		pkgPath = named.Obj().Pkg().Path() + "."
	}
	k := pkgPath + named.Obj().Name()
	if isPointer {
		k = "*" + k
	}

	if v, ok := d.visited[k].(*Named); ok {
		return v, nil
	}

	pkg, err := d.normalizePkg(named.Obj().Pkg())
	if err != nil {
		return nil, err
	}

	title, description, tags, err := d.commentsAndTagsFind(named.Obj().Name(), named.Obj().Pos())
	if err != nil {
		return nil, err
	}

	nt = &Named{
		origin:      named,
		Pkg:         pkg,
		Title:       title,
		Description: description,
		Name:        named.Obj().Name(),
		IsPointer:   isPointer,
		IsExported:  named.Obj().Exported(),
		IsAlias:     named.Obj().IsAlias(),
		Position:    d.pkg.Fset.Position(named.Obj().Pos()),
		Tags:        tags,
	}

	d.visited[k] = nt

	nt.Type, err = d.normalizeRecursive(named.Obj().Type().Underlying(), false)
	if err != nil {
		return nil, err
	}

	_, isStruct := nt.Type.(*Struct)

	for i := 0; i < named.NumMethods(); i++ {
		method, err := d.normalizeFunc(named.Method(i))
		if err != nil {
			return nil, err
		}
		if isStruct {
			returns, err := d.findFuncReturn(method)
			if err != nil {
				return nil, err
			}
			method.Returns = returns
		}
		nt.Methods = append(nt.Methods, method)
	}
	return nt, nil
}

func (d *Decoder) normalizeVar(t *stdtypes.Var) (*Var, error) {
	varType, err := d.normalizeRecursive(t.Type(), false)
	if err != nil {
		return nil, err
	}
	title, description, tags, err := d.commentsAndTagsFind(t.Name(), t.Pos())
	if err != nil {
		return nil, err
	}
	return &Var{
		Name:      t.Name(),
		Embedded:  t.Embedded(),
		Exported:  t.Exported(),
		IsField:   t.IsField(),
		IsContext: IsContext(varType),
		IsError:   IsError(varType),
		Type:      varType,
		Zero:      zeroValue(t.Type().Underlying()),
		Title:     title + "\n" + description,
		Tags:      tags,
		Position:  d.pkg.Fset.Position(t.Pos()),
	}, nil
}

func (d *Decoder) normalizeSignature(t *stdtypes.Signature) (st *Sign, err error) {
	st = &Sign{
		IsVariadic: t.Variadic(),
	}
	if t.Recv() != nil {
		st.Recv, err = d.normalizeRecursive(t.Recv().Type(), false)
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < t.Params().Len(); i++ {
		v := t.Params().At(i)
		nv, err := d.normalizeVar(v)
		if err != nil {
			return nil, err
		}
		st.Params = append(st.Params, nv)
	}
	if t.Variadic() {
		st.Params[len(st.Params)-1].IsVariadic = true
	}
	for i := 0; i < t.Results().Len(); i++ {
		v := t.Results().At(i)
		nv, err := d.normalizeVar(v)
		if err != nil {
			return nil, err
		}
		if nv.Name != "" {
			st.IsNamed = true
		}
		st.Results = append(st.Results, nv)
	}
	return st, nil
}

func (d *Decoder) normalizeFunc(t *stdtypes.Func) (*Func, error) {
	title, description, tags, err := d.commentsAndTagsFind(t.Name(), t.Pos())
	if err != nil {
		return nil, err
	}
	fnSig := t.Type().(*stdtypes.Signature)
	sig, err := d.normalizeSignature(fnSig)
	if err != nil {
		return nil, err
	}
	pkg, err := d.normalizePkg(t.Pkg())
	if err != nil {
		return nil, err
	}
	fn := &Func{
		Pkg:         pkg,
		FullName:    t.FullName(),
		Name:        t.Name(),
		Exported:    t.Exported(),
		Sig:         sig,
		Title:       title,
		Description: description,
		Tags:        tags,
		Position:    d.pkg.Fset.Position(t.Pos()),
		Text:        stdtypes.NewFunc(0, nil, t.Name(), stdtypes.NewSignatureType(nil, nil, nil, fnSig.Params(), fnSig.Results(), fnSig.Variadic())).String(),
	}
	return fn, nil
}

func (d *Decoder) normalizeRecursive(t any, isPointer bool) (any, error) {
	switch t := t.(type) {
	case *packages.Module:
		return d.normalizeModule(t)
	case *stdtypes.PkgName:
		return d.normalizePkg(t.Pkg())
	case *stdtypes.Var:
		return d.normalizeVar(t)
	case *stdtypes.Func:
		return d.normalizeFunc(t)
	case *stdtypes.Map:
		return d.normalizeMap(t.Key(), t.Elem(), isPointer)
	case *stdtypes.Slice:
		return d.normalizeSlice(t.Elem(), isPointer)
	case *stdtypes.Array:
		return d.normalizeArray(t.Elem(), t.Len(), isPointer)
	case *stdtypes.Pointer:
		return d.normalizeRecursive(t.Elem(), true)
	case *stdtypes.Struct:
		return d.normalizeStruct(t, isPointer)
	case *stdtypes.Signature:
		return d.normalizeSignature(t)
	case *stdtypes.Interface:
		return d.normalizeInterface(t)
	case *stdtypes.Named:
		return d.normalizeNamed(t, isPointer)
	case *stdtypes.TypeName:
		return d.normalizeRecursive(t.Type(), isPointer)
	case *stdtypes.Basic:
		return d.normalizeBasic(t, isPointer)
	case *stdtypes.Chan:
		return d.normalizeChan(t)
	case *stdtypes.Tuple:
		return d.normalizeTuple(t)
	}
	return nil, fmt.Errorf("unknown type: %T", t)
}

func (d *Decoder) findFuncReturn(targetFn *Func) (values []any, err error) {
	err = traverseDecls(d.packages, func(pkg *packages.Package, file *ast.File, decl ast.Decl) error {
		if t, ok := decl.(*ast.FuncDecl); ok {
			obj := pkg.TypesInfo.ObjectOf(t.Name)
			if obj == nil {
				return nil
			}

			fn, ok := obj.(*stdtypes.Func)
			if !ok {
				return nil
			}
			if !fn.Exported() {
				return nil
			}
			if targetFn.FullName != fn.FullName() {
				return nil
			}

			sig := pkg.TypesInfo.TypeOf(t.Name).(*stdtypes.Signature)

			recv := sig.Recv()
			if recv == nil {
				return nil
			}
			values, err = d.loadResultsRecursive(pkg, t.Body.List)
			if err != nil {
				return err
			}

		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

func (d *Decoder) loadResultsRecursive(pkg *packages.Package, stmts []ast.Stmt) (results []any, err error) {
	for _, stmt := range stmts {
		switch tp := stmt.(type) {
		case *ast.IfStmt:
			items, err := d.loadResultsRecursive(pkg, tp.Body.List)
			if err != nil {
				return nil, err
			}
			results = append(results, items...)
		case *ast.ReturnStmt:
			for _, result := range tp.Results {
				if callExpr, ok := result.(*ast.CallExpr); ok {
					if fSel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if t, ok := pkg.TypesInfo.TypeOf(fSel.X).(*stdtypes.Named); ok {
							n, err := d.normalizeNamed(t, true)
							if err != nil {
								return nil, err
							}
							results = append(results, &Call{FuncName: fSel.Sel.Name, Named: n})
							continue
						}
					}
				}

				if v, ok := pkg.TypesInfo.Types[result]; ok {
					tv := constant.Val(v.Value)
					if tv != nil {
						results = append(results, tv)
					} else {
						entity, err := d.normalizeRecursive(v.Type, false)
						if err != nil {
							return nil, err
						}
						results = append(results, entity)
					}
				}
			}
		}
	}
	return
}

func (d *Decoder) commentsAndTagsFind(name string, pos token.Pos) (title, description string, tags Tags, err error) {
	var tagComments Comments
	allComments := d.commentFind(name, pos)
	for _, comment := range allComments {
		if comment.IsTag {
			tagComments = append(tagComments, comment)
		} else if comment.IsTitle {
			title = comment.Value
		} else {
			description += comment.Value + "\n"
		}
	}
	if len(tagComments) > 0 {
		tags, err = parseTags(tagComments)
		if err != nil {
			return
		}
	}
	return
}

func (d *Decoder) commentFind(name string, pos token.Pos) (result Comments) {
	position := d.pkg.Fset.Position(pos)
	for _, pkg := range d.packages {
		for _, file := range pkg.Syntax {
			for _, commentGroup := range file.Comments {
				cg := d.pkg.Fset.Position(commentGroup.End())
				if cg.Line == position.Line-1 && cg.Filename == position.Filename {
					for _, comment := range commentGroup.List {
						text := strings.TrimLeft(strings.TrimLeft(comment.Text, "//"), " ")
						isTitle := strings.HasPrefix(text, name)
						isTag := strings.HasPrefix(text, "@")
						if isTitle {
							text = strings.Replace(text, name+" ", "", -1)
						}
						result = append(result, &Comment{
							Value:    text,
							IsTitle:  isTitle,
							IsTag:    isTag,
							Position: d.pkg.Fset.Position(comment.End()),
						})
					}
				}
			}
		}
	}
	return
}

func NewDecoder(pkg *packages.Package, packages []*stdpackages.Package) *Decoder {
	return &Decoder{pkg: pkg, packages: packages, visited: make(map[string]any, 1024)}
}
