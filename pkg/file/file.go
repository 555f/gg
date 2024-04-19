package file

type File interface {
	Path() string
	Bytes() ([]byte, error)
}
