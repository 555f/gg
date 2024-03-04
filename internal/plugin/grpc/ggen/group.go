package pgen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

type Group struct {
	name      string
	items     []Code
	open      string
	close     string
	separator string
	multi     bool
}

func (g *Group) render(f *File, w io.Writer, s *Statement) error {
	if g.open != "" {
		if _, err := w.Write([]byte(g.open)); err != nil {
			return err
		}
	}
	isNull, err := g.renderItems(f, w)
	if err != nil {
		return err
	}
	if !isNull && g.multi && g.close != "" {
		s := "\n"
		if g.separator == "," {
			s = ",\n"
		}
		if _, err := w.Write([]byte(s)); err != nil {
			return err
		}
	}
	if g.close != "" {
		if _, err := w.Write([]byte(g.close)); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) renderItems(f *File, w io.Writer) (isNull bool, err error) {
	if len(g.items) == 0 {
		return false, nil
	}

	first := true
	for _, code := range g.items {
		if code == nil {
			continue
		}
		if pt, ok := code.(token); ok && pt.typ == packageToken {
			f.register(pt.content.(string))
		}
		if !first && g.separator != "" {
			if _, err := w.Write([]byte(g.separator)); err != nil {
				return false, err
			}
		}
		if g.multi {
			if _, err := w.Write([]byte("\n")); err != nil {
				return false, err
			}
		}
		if err := code.render(f, w, nil); err != nil {
			return false, err
		}
		first = false
	}
	if g.separator != "" {
		if _, err := w.Write([]byte(g.separator)); err != nil {
			return false, err
		}
	}
	return first, nil
}

func (g *Group) Render(writer io.Writer) error {
	return g.RenderWithFile(writer, NewFile(""))
}

func (g *Group) String() string {
	buf := bytes.Buffer{}
	if err := g.Render(&buf); err != nil {
		panic(err)
	}
	return buf.String()
}

func (g *Group) RenderWithFile(writer io.Writer, file *File) error {
	buf := &bytes.Buffer{}
	if err := g.render(file, buf, nil); err != nil {
		return err
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error %s while formatting source:\n%s", err, buf.String())
	}
	if _, err := writer.Write(b); err != nil {
		return err
	}
	return nil
}
