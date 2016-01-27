package config

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	f       = new(File)
	yamlCfg *Config
	jsonCfg *Config
)

func TestFileFromArg(t *testing.T) {
	loadFromArgs(t, "yaml")
	cleanArgs()
	loadFromArgs(t, "json")
}

func TestJsonRemoveComments(t *testing.T) {
	f.Data = []byte(`{
    // hey
    // wat
  }`)

	f.RemoveJSONComments()
	assert.Equal(t, `{
    
    
  }`, string(f.Data))
}

func TestJsonLoad(t *testing.T) {
	var err error
	f.FilePath, err = os.Getwd()
	assert.NoError(t, err)
	assert.NotEmpty(t, f.FilePath)
	f.FilePath += "/../_example/simple/b0x.json"

	cfg, err := f.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.NotEqual(t, `{
    
    
  }`, string(f.Data))

	cfg.Defaults()
	jsonCfg = cfg
}

func TestYamlLoad(t *testing.T) {
	loadFromArgs(t, "yaml")

	var err error
	f.FilePath, err = os.Getwd()
	assert.NoError(t, err)
	assert.NotEmpty(t, f.FilePath)
	f.FilePath += "/../_example/simple/b0x.yaml"

	cfg, err := f.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, string(f.Data))

	cfg.Defaults()
	yamlCfg = cfg
}

func TestComparison(t *testing.T) {
	assert.True(t, reflect.DeepEqual(yamlCfg, jsonCfg))
}

func loadFromArgs(t *testing.T, ext string) {
	// insert "b0x.ext" to args
	os.Args = append(os.Args, "b0x."+ext)

	// test b0x.ext from args
	err := f.FromArg()
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(f.FilePath, "b0x."+ext))
}

func cleanArgs() {
	// remove "b0x.*" from last arg
	os.Args = os.Args[:len(os.Args)-1]
}
