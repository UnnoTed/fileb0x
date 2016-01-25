package file

// File holds file's data
type File struct {
	Name  string
	Path  string
	Data  string
	Bytes []byte
}

// NewFile creates a new File
func NewFile() *File {
	f := new(File)
	return f
}
