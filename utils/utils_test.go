package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixPath(t *testing.T) {
	p := `C:\test\a\b\c\d\e\f`
	fixedPath := FixPath(p)

	assert.Equal(t, `C:/test/a/b/c/d/e/f`, fixedPath)
}

func TestFixName(t *testing.T) {
	name := `notporn/empty folder/ufo, porno.flv`
	fixedName := FixName(name)

	assert.Equal(t, `notporn_empty-folder_ufo__-porno.flv`, fixedName)
}

// https://github.com/UnnoTed/fileb0x/issues/8
func TestIssue8(t *testing.T) {
	type replacer struct {
		file   string
		base   string
		prefix string
		result string
	}

	list := []replacer{
		{"./main.go", "./", "_bug?_", "_bug?_main.go"},
		{"./fileb0x.test.yaml", "./", "_bug?_", "_bug?_fileb0x.test.yaml"},
		{"./static/test.txt", "", "test_prefix/", "test_prefix/static/test.txt"},
		{"./static/test.txt", "./static/", "", "test.txt"},
	}

	fixit := func(r replacer) string {
		r.file = FixPath(r.file)

		var fixedPath string
		if r.prefix != "" || r.base != "" {
			if strings.HasPrefix(r.base, "./") {
				r.base = r.base[2:]
			}

			if strings.HasPrefix(r.file, r.base) {
				fixedPath = r.prefix + r.file[len(r.base):]
			} else {
				if r.base != "" {
					fixedPath = r.prefix + r.file
				}
			}

			fixedPath = FixPath(fixedPath)
		} else {
			fixedPath = FixPath(r.file)
		}

		return fixedPath
	}

	for _, r := range list {
		assert.Equal(t, r.result, fixit(r))
	}
}

func TestGetCurrentDir(t *testing.T) {
	dir, err := GetCurrentDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)

	assert.Contains(t, dir, "fileb0x")
	assert.Contains(t, dir, "utils")
}

func TestExists(t *testing.T) {
	dir, err := GetCurrentDir()
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)

	exists := Exists(dir + "/_testmain.go")
	assert.True(t, exists)

	exists = Exists("wat")
	assert.False(t, exists)
}
