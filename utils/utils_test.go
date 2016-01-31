package utils

import (
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
