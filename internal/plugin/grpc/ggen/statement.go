package pgen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

type Statement []Code

func newStatement() *Statement {
	return &Statement{}
}

func (s *Statement) Clone() *Statement {
	return &Statement{s}
}

func (s *Statement) render(f *File, w io.Writer, _ *Statement) error {
	first := true
	for _, code := range *s {
		if code == nil {
			continue
		}
		if !first {
			if _, err := w.Write([]byte(" ")); err != nil {
				return err
			}
		}
		if err := code.render(f, w, s); err != nil {
			return err
		}
		first = false
	}
	return nil
}

func (s *Statement) Render(writer io.Writer) error {
	return s.RenderWithFile(writer, NewFile(""))
}

func (s *Statement) String() string {
	buf := bytes.Buffer{}
	if err := s.Render(&buf); err != nil {
		panic(err)
	}
	return buf.String()
}

func (s *Statement) RenderWithFile(writer io.Writer, file *File) error {
	buf := &bytes.Buffer{}
	if err := s.render(file, buf, nil); err != nil {
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
