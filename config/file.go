package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/UnnoTed/fileb0x/utils"

	"gopkg.in/yaml.v2"
)

// File holds config file info
type File struct {
	FilePath string
	Data     []byte
	Mode     string // "json" || "yaml" || "yml"
}

// FromArg gets the json/yaml file from args
func (f *File) FromArg(read bool) error {
	// (length - 1)
	arg := os.Args[len(os.Args)-1:][0]

	// get extension
	ext := path.Ext(arg)
	if len(ext) > 1 {
		ext = ext[1:] // remove dot
	}

	// when json/yaml file isn't found on last arg
	// it searches for a ".json", ".yaml" or ".yml" string in all args
	if ext != "json" && ext != "yaml" && ext != "yml" {
		// loop through args
		for _, a := range os.Args {
			// get extension
			ext := path.Ext(a)

			// check for valid extensions
			if ext == ".json" || ext == ".yaml" || ext == ".yml" {
				f.Mode = ext[1:] // remove dot
				ext = f.Mode
				arg = a
				break
			}
		}
	} else {
		f.Mode = ext
	}

	// check if extension is json or yaml
	// then get it's absolute path
	if ext == "json" || ext == "yaml" || ext == "yml" {
		abs := filepath.IsAbs(arg)
		if !abs {
			dir, err := utils.GetCurrentDir()
			if err != nil {
				return err
			}

			arg = filepath.Clean(dir + "/" + arg)
		}

		f.FilePath = arg

		// so we can test without reading a file
		if read {
			if !utils.Exists(f.FilePath) {
				return errors.New("Error: I Can't find the config file at [" + f.FilePath + "]")
			}
		}
	} else {
		return errors.New("Error: You must specify a json or yaml file")
	}

	return nil
}

// Parse gets the config file's content from File.Data
func (f *File) Parse() (*Config, error) {
	// remove comments
	f.RemoveJSONComments()

	var to *Config
	var err error

	// parse file
	if f.Mode == "json" {
		err = json.Unmarshal(f.Data, &to)
	} else if f.Mode == "yaml" || f.Mode == "yml" {
		err = yaml.Unmarshal(f.Data, &to)
	}

	return to, err
}

// Load the json/yaml file that was specified from args
// and transform it into a config struct
func (f *File) Load() (*Config, error) {
	var err error
	if !utils.Exists(f.FilePath) {
		return nil, errors.New("Error: I Can't find the config file at [" + f.FilePath + "]")
	}

	// read file
	f.Data, err = ioutil.ReadFile(f.FilePath)
	if err != nil {
		return nil, err
	}

	// parse file
	return f.Parse()
}

// RemoveJSONComments from the file
func (f *File) RemoveJSONComments() {
	if f.Mode == "json" {
		// remove inline comments
		f.Data = []byte(regexComments.ReplaceAllString(string(f.Data), ""))
	}
}
