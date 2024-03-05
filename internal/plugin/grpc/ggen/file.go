package pgen

import (
	"bytes"
)

func NewFile(packageName string) *File {
	return &File{
		Group: &Group{
			multi: true,
		},
		name:    packageName,
		imports: map[string]struct{}{},
	}
}

type File struct {
	*Group
	name     string
	syntax   string
	imports  map[string]struct{}
	comments []string
}

func (f *File) Syntax(syntax string) {
	f.syntax = syntax
}

func (f *File) Comment(comment string) {
	f.comments = append(f.comments, comment)
}

func (f *File) ImportName(name string) {
	f.imports[name] = struct{}{}
}

func (f *File) ImportNames(names map[string]string) {
	for name := range names {
		f.imports[name] = struct{}{}
	}
}

func (f *File) register(name string) {
	if path := standardLibrary[name]; path != "" {
		f.imports[path] = struct{}{}
	}
}

func (f *File) String() string {
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		panic(err)
	}
	return buf.String()
}
