package file

import (
	"bytes"
	"fmt"
)

var _ File = &TxtFile{}

type TxtFile struct {
	buf  bytes.Buffer
	path string
}

func (f *TxtFile) Line() {
	f.WriteText("\n")
}

func (f *TxtFile) WriteBytes(p []byte) {
	_, _ = f.buf.Write(p)
}

func (f *TxtFile) Write(p []byte) (n int, err error) {
	return f.buf.Write(p)
}

func (f *TxtFile) WriteText(format string, a ...any) {
	_, _ = fmt.Fprintf(&f.buf, format, a...)
}

func (f *TxtFile) Filepath() string {
	return f.path
}

func (f *TxtFile) Bytes() ([]byte, error) {
	return f.buf.Bytes(), nil
}

func NewTxtFile(path string) *TxtFile {
	return &TxtFile{path: path}
}
