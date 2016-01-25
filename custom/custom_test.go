package custom

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"

	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/stretchr/testify/assert"
)

func TestCustomParse(t *testing.T) {
	c := new(Custom)
	c.Files = []string{
		"../_example/simple/public/",
	}

	c.Base = "../_example/simple/"
	c.Prefix = "prefix_test/"
	c.Exclude = []string{
		"public/assets/data/exclude_me.txt",
	}

	c.Replace = []Replacer{
		{
			File: "public/assets/data/test*.json",
			Replace: map[string]string{
				"{world}": "earth",
				"{EMAIL}": "aliens@nasa.com",
			},
		},
	}

	files := make(map[string]*file.File)
	dirs := new(dir.Dir)

	err := c.Parse(&files, &dirs)
	assert.NoError(t, err)
	assert.NotNil(t, files)
	assert.NotNil(t, dirs)

	for _, f := range files {
		assert.True(t, strings.HasPrefix(f.Path, c.Prefix))
		assert.NotEqual(t, "exclude_me.txt", f.Name)

		if f.Name == "test1.json" {
			e := "{\r\n  \"he\": \"llo\",\r\n  \"replace_test\": \"earth\"\r\n}"
			assert.Equal(t, e, data2str(f.Data))

		} else if f.Name == "test2.json" {
			e := "{\r\n  \"email\": \"aliens@nasa.com\"\r\n}"
			assert.Equal(t, e, data2str(f.Data))
		}
	}

	ds := dirs.Clean()
	var blacklist []string
	for _, d := range ds {
		assert.True(t, strings.HasPrefix(d, c.Prefix))
		assert.NotContains(t, blacklist, d)
		blacklist = append(blacklist, d)
	}
}

func data2str(h string) string {
	h = strings.TrimPrefix(h, `[]byte("`)
	h = strings.TrimSuffix(h, `")`)
	h = strings.Replace(h, `\x`, "", -1)

	b, err := hex.DecodeString(h)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
