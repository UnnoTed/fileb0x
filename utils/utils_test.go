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
