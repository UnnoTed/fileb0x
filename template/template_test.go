package template

import (
	"strings"
	"testing"

	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	files := make(map[string]*file.File)
	files["test_file.txt"] = &file.File{
		Name: "test_file.txt",
		Path: "static/test_file.txt",
		Data: `[]byte("\x12\x34\x56\x78\x10")`,
	}

	dirs := new(dir.Dir)
	dirs.Insert("static/")

	tp := new(Template)
	tp.Set("files")
	tp.Variables = struct {
		Pkg     string
		Files   map[string]*file.File
		Spread  bool
		DirList []string
	}{
		Pkg:     "main",
		Files:   files,
		Spread:  false,
		DirList: dirs.Clean(),
	}

	tmpl, err := tp.Exec()
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpl)

	s := string(tmpl)

	assert.True(t, strings.Contains(s, `var Filestatictestfiletxt = []byte("\x12\x34\x56\x78\x10")`))
	assert.True(t, strings.Contains(s, `err = FS.Mkdir("static/", 0777)`))
	assert.True(t, strings.Contains(s, `f, err = FS.OpenFile("static/test_file.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)`))

	// now with spread
	tp.Set("file")
	tp.Variables = struct {
		Pkg  string
		Path string
		Name string
		Dir  [][]string
		Data string
	}{
		Pkg:  "main",
		Path: files["test_file.txt"].Path,
		Name: files["test_file.txt"].Name,
		Dir:  dirs.List,
		Data: files["test_file.txt"].Data,
	}

	tmpl, err = tp.Exec()
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpl)

	s = string(tmpl)

	assert.True(t, strings.Contains(s, `var Filestatictestfiletxt = []byte("\x12\x34\x56\x78\x10")`))
	assert.True(t, strings.Contains(s, `f, err := FS.OpenFile("static/test_file.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)`))
}
