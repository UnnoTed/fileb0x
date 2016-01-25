package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {
	cfg := new(Config)
	assert.Equal(t, "", cfg.Dest)
	assert.Equal(t, "", cfg.Output)
	assert.Equal(t, "", cfg.Pkg)

	err := cfg.Defaults()
	assert.NoError(t, err)

	assert.Equal(t, "/", cfg.Dest)
	assert.Equal(t, "ab0x.go", cfg.Output)
	assert.Equal(t, "main", cfg.Pkg)
}
