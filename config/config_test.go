package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	var err error
	cfgList := []*Config{
		{Dest: "", Pkg: "", Output: ""},
		{Dest: "hello", Pkg: "static", Output: "ssets"},
		{Dest: "hello", Pkg: "static", Output: "amissingEnd"},
		{Dest: "hello", Pkg: "static", Output: "compiled", NoPrefix: true},
	}

	expecTedList := []*Config{
		{Dest: "/", Pkg: "main", Output: "ab0x.go"},
		{Dest: "hello/", Pkg: "static", Output: "assets.go"},
		{Dest: "hello/", Pkg: "static", Output: "amissingEnd.go"},
		{Dest: "hello/", Pkg: "static", Output: "compiled.go", NoPrefix: true},
	}

	for i, cfg := range cfgList {
		err = cfg.Defaults()
		assert.NoError(t, err)

		err = expecTedList[i].Defaults()
		assert.NoError(t, err)

		eq := reflect.DeepEqual(cfg, expecTedList[i])
		assert.True(t, eq, "NOT EQUAL:", cfg, expecTedList[i])
	}
}
