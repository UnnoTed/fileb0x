package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/UnnoTed/fileb0x/compression"
	"github.com/UnnoTed/fileb0x/custom"
	"github.com/UnnoTed/fileb0x/updater"
)

// Config holds the json/yaml/toml data
type Config struct {
	Dest        string
	Pkg         string
	Fmt         bool // gofmt
	Compression *compression.Options
	Tags        string

	Output string

	Custom []custom.Custom

	Spread     bool
	Unexported bool
	Clean      bool
	Debug      bool
	Updater    updater.Config
}

// Defaults set the default value for some variables
func (cfg *Config) Defaults() error {
	// default destination
	if cfg.Dest == "" {
		cfg.Dest = "/"
	}

	// insert "/" at end of dest when it's not found
	if !strings.HasSuffix(cfg.Dest, "/") {
		cfg.Dest += "/"
	}

	// default file name
	if cfg.Output == "" {
		cfg.Output = "b0x.go"
	}

	// inserts .go at the end of file name
	if !strings.HasSuffix(cfg.Output, ".go") {
		cfg.Output += ".go"
	}

	// inserts an A before the output file's name so it can
	// run init() before b0xfile's
	if !strings.HasPrefix(cfg.Output, "a") {
		cfg.Output = "a" + cfg.Output
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

		// remove matched file
		for _, f := range matches {
			err = os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
