package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	j = new(JSON)
)

func TestJsonFileFromArg(t *testing.T) {
	err := j.FromArg()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(j.FilePath, "b0x.json"))
}

func TestJsonRemoveComments(t *testing.T) {
	j.Data = []byte(`{
    // hey
    // wat
  }`)

	j.RemoveComments()
	assert.Equal(t, j.Data, []byte(`{
    
    
  }`))
}

func TestJsonLoad(t *testing.T) {
	var err error
	j.FilePath, err = os.Getwd()
	assert.NoError(t, err)
	assert.NotEmpty(t, j.FilePath)
	j.FilePath += "/../_example/simple/b0x.json"

	cfg, err := j.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.NotEqual(t, j.Data, []byte(`{
    
    
  }`))
}
