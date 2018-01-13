package template

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugRemapPath(t *testing.T) {
	remap := map[string]map[string]string{
		"public/README.md": {
			"prefix": "public/",
			"base":   "../../",
		},
	}

	open := func(path string) string {
		path = strings.TrimPrefix(path, "/")

		for current, f := range remap {
			if path == current {
				path = f["base"] + strings.TrimPrefix(path, f["prefix"])
				break
			}
		}

		return path
	}

	tests := map[string]string{
		"/public/README.md": "../../README.md",
	}

	for actual, expected := range tests {
		assert.Equal(t, expected, open(actual))
		t.Log(expected, actual, open(actual))
	}
}
