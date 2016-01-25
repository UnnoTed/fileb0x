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

	cfg.Defaults()

	assert.Equal(t, "/", cfg.Dest)
	assert.Equal(t, "b0x.go", cfg.Output)
	assert.Equal(t, "main", cfg.Pkg)
}
