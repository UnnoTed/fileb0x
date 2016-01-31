package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/UnnoTed/fileb0x/utils"
	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	var err error
	cfgList := []*Config{
		{Dest: "", Pkg: "", Output: ""},
		{Dest: "hello", Pkg: "static", Output: "ssets"},
		{Dest: "hello", Pkg: "static", Output: "amissingEnd"},
	}

	expecTedList := []*Config{
		{Dest: "/", Pkg: "main", Output: "ab0x.go"},
		{Dest: "hello/", Pkg: "static", Output: "assets.go"},
		{Dest: "hello/", Pkg: "static", Output: "amissingEnd.go"},
	}

	for i, cfg := range cfgList {
		err = cfg.Defaults()
		assert.NoError(t, err)

		eq := reflect.DeepEqual(cfg, expecTedList[i])
		assert.True(t, eq, "NOT EQUAL:", cfg, expecTedList[i])
	}
}

func TestConfigClean(t *testing.T) {
	var err error
	cfg := new(Config)
	cfg.Dest = "./clean_test123456789/"
	cfg.Clean = true

	err = os.Mkdir(cfg.Dest, 0777)
	assert.NoError(t, err)

	var files = []string{
		"a.go",
		"b.go",
		"c.go",
		"d.go",
		"e.go",
		"f.go",
		"g.go",
		"h.go",
	}

	// write files into the folder so it can be cleaned by Config.Defaults()
	nothing := []byte("nothing")
	for _, file := range files {
		err = ioutil.WriteFile(cfg.Dest+"b0xfile_"+file, nothing, 0777)
		assert.NoError(t, err)
		assert.True(t, utils.Exists(cfg.Dest+"b0xfile_"+file))
	}

	// clean
	err = cfg.Defaults()
	assert.NoError(t, err)

	// verify that the files were deleTed
	for _, file := range files {
		assert.False(t, utils.Exists(cfg.Dest+"b0xfile_"+file))
	}

	// delete dest folder
	os.RemoveAll(cfg.Dest)
	assert.False(t, utils.Exists(cfg.Dest))
}
