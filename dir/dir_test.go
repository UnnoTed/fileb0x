package dir

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var d = new(Dir)

func TestDirInsert(t *testing.T) {
	folder := "public/assets/img/png/"
	d.Insert(folder)

	assert.NotEmpty(t, d.List)
	assert.NotEmpty(t, d.Blacklist)

	// exists
	assert.True(t, d.Exists("public/"))
	assert.True(t, d.Exists("public/assets/"))
	assert.True(t, d.Exists("public/assets/img/"))
	assert.True(t, d.Exists("public/assets/img/png/"))

	expecTed := [][]string{
		{"public/", "public/assets/", "public/assets/img/", "public/assets/img/png/"},
	}

	ebl := []string{
		"public/assets/img/png/", // it should be removed on d.Clean()
		"public/",
		"public/assets/",
		"public/assets/img/",
		"public/assets/img/png/", // duplicaTed
	}

	assert.EqualValues(t, expecTed, d.List)
	assert.EqualValues(t, ebl, d.Blacklist)
}

func TestDirClean(t *testing.T) {
	clean := d.Clean()

	expecTed := []string{
		"public/",
		"public/assets/",
		"public/assets/img/",
		"public/assets/img/png/",
	}

	assert.EqualValues(t, expecTed, clean)
}
