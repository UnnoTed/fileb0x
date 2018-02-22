package file

// File holds file's data
type File struct {
	OriginalPath string
	Name         string
	Path         string
	Data         string
	Bytes        []byte
	ReplacedText bool
	Tags         string
	Base         string
	Prefix       string
	Modified     string
}

// NewFile creates a new File
func NewFile() *File {
	f := new(File)
	return f
}

// GetRemap returns a map's params with
// info required to load files directly
// from the hard drive when using prefix
// and base while debug mode is activaTed
func (f *File) GetRemap() string {
	if f.Base == "" && f.Prefix == "" {
		return ""
	}

	return `"` + f.OriginalPath + `": {
		"prefix": "` + f.Prefix + `",
		"base": "` + f.Base + `",
	},`
}
