package config

import (
	"os"
	"path/filepath"

	"github.com/UnnoTed/fileb0x/custom"
)

// Config holds the json data
type Config struct {
	Dest string
	Pkg  string

	Output string

	Custom []custom.Custom

	Spread     bool
	Unexported bool
	Clean      bool
}

// Defaults set the default value for some variables
func (cfg *Config) Defaults() error {
	// default destination
	if cfg.Dest == "" {
		cfg.Dest = "/"
	}

	// insert "/" at end of dest when it's not found
	if cfg.Dest[len(cfg.Dest)-1:] != "/" {
		cfg.Dest = "/"
	}

	// default file name
	if cfg.Output == "" {
		cfg.Output = "b0x.go"
	}

	// default package
	if cfg.Pkg == "" {
		cfg.Pkg = "main"
	}

	// remove b0xfiles when [clean] is true
	// it doesn't clean destination's folders
	if cfg.Clean {
		matches, err := filepath.Glob(cfg.Dest + "b0xfile_*.go")
		if err != nil {
			return err
		}

		for _, f := range matches {
			err = os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
