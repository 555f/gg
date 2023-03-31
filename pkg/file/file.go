package file

type File interface {
	Filepath() string
	Bytes() ([]byte, error)
}
