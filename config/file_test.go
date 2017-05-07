package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	f        = new(File)
	yamlCfg  *Config
	jsonCfg  *Config
	confYAML []byte
	confJSON []byte

	yamlPath = "/../_example/simple/b0x.yaml"
	jsonPath = "/../_example/simple/b0x.json"
)

func TestFileFromArg(t *testing.T) {
	loadFromArgs(t, "yaml")

	// remove "b0x.*" from last arg
	os.Args = os.Args[:len(os.Args)-1]

	err := f.FromArg(false)
	assert.Error(t, err)
	assert.Equal(t, "Error: You must specify a json, yaml or toml file", err.Error())

	err = f.FromArg(true)
	assert.Error(t, err)
	assert.Equal(t, "Error: You must specify a json, yaml or toml file", err.Error())

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
	f.FilePath += jsonPath

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

	cfg, err := f.Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Empty(t, string(f.Data))

	f.FilePath += yamlPath
	cfg, err = f.Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, string(f.Data))

	cfg.Defaults()
	yamlCfg = cfg
}

func TestYamlParse(t *testing.T) {
	var err error
	f.Data = nil
	f.Mode = "yaml"

	f.Data, err = ioutil.ReadFile("." + yamlPath)
	assert.NoError(t, err)
	assert.NotNil(t, f.Data)

	cfg, err := f.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestJsonParse(t *testing.T) {
	var err error
	f.Data = nil
	f.Mode = "json"

	f.Data, err = ioutil.ReadFile("." + jsonPath)
	assert.NoError(t, err)
	assert.NotNil(t, f.Data)

	cfg, err := f.Parse()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestComparison(t *testing.T) {
	assert.True(t, reflect.DeepEqual(yamlCfg, jsonCfg))
}

func loadFromArgs(t *testing.T, ext string) {
	// insert "b0x.ext" to args
	os.Args = append(os.Args, "b0x."+ext)

	err := f.FromArg(false)
	assert.NoError(t, err)

	err = f.FromArg(true)
	assert.Error(t, err)
	assert.Equal(t, "Error: I Can't find the config file at ["+f.FilePath+"]", err.Error())

	// test b0x.ext from args
	err = f.FromArg(false)
	assert.NoError(t, err)
	assert.True(t, strings.HasSuffix(f.FilePath, "b0x."+ext))
}
