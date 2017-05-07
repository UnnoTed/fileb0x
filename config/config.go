package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/UnnoTed/fileb0x/compression"
	"github.com/UnnoTed/fileb0x/custom"
)

type Updater struct {
	Username string
	Password string
	Enabled  bool
	Port     int
}

func (u Updater) CheckInfo() error {
	if !u.Enabled {
		return nil
	}

	if u.Username == "{FROM_ENV}" || u.Username == "" {
		u.Username = os.Getenv("fileb0x_username")
	}

	if u.Password == "{FROM_ENV}" || u.Password == "" {
		u.Password = os.Getenv("fileb0x_password")
	}

	// check for empty username and password
	if u.Username == "" {
		return errors.New("fileb0x: You must provide an username in the config file or through an env var: fileb0x_username")

	} else if u.Password == "" {
		return errors.New("fileb0x: You must provide an password in the config file or through an env var: fileb0x_password")
	}

	return nil
}

// Config holds the json/yaml/toml data
type Config struct {
	Dest        string
	Pkg         string
	Fmt         bool // gofmt
	Compression *compression.Options

	Output string

	Custom []custom.Custom

	Spread     bool
	Unexported bool
	Clean      bool
	Debug      bool
	Updater    Updater
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
