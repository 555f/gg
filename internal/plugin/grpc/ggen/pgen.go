package pgen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type Code interface {
	render(f *File, w io.Writer, s *Statement) error
}

func (f *File) Save(filename string) error {
	// notest
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func (f *File) Render(w io.Writer) error {
	body := &bytes.Buffer{}
	if err := f.render(f, body, nil); err != nil {
		return err
	}
	source := &bytes.Buffer{}
	if f.syntax != "" {
		f.renderSyntax(source, f.syntax)
	}
	if len(f.comments) > 0 {
		for _, c := range f.comments {
			if err := Comment(c).render(f, source, nil); err != nil {
				return err
			}
			if _, err := fmt.Fprint(source, "\n"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(source, "\n"); err != nil {
			return err
		}
	}
	if err := f.renderImports(source); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(source, "package %s;", f.name); err != nil {
		return err
	}
	if _, err := fmt.Fprint(source, "\n\n"); err != nil {
		return err
	}
	if _, err := source.Write(body.Bytes()); err != nil {
		return err
	}
	output := source.Bytes()
	if _, err := w.Write(output); err != nil {
		return err
	}
	return nil
}

func (f *File) renderImports(source io.Writer) error {
	if len(f.imports) == 0 {
		return nil
	}
	paths := []string{}
	for path := range f.imports {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if _, err := fmt.Fprintf(source, "import %s;\n", strconv.Quote(path)); err != nil {
			return err
		}
	}
	fmt.Fprintf(source, "\n")
	return nil
}

func (f *File) renderSyntax(source io.Writer, syntax string) error {
	if _, err := fmt.Fprintf(source, "syntax = %s;\n\n", strconv.Quote(syntax)); err != nil {
		return err
	}
	return nil
}
