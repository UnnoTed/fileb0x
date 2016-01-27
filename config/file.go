package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// File holds config file info
type File struct {
	FilePath string
	Data     []byte
	Mode     string // "json" || "yaml"
}

// FromArg gets the json/yaml file from args
func (f *File) FromArg() error {
	// (length - 1)
	arg := os.Args[len(os.Args)-1:][0]

	// get extension
	ext := path.Ext(arg)
	ext = ext[1:] // remove dot

	// when json/yaml file isn't found on last arg
	// it searches for a ".json" or ".yaml" string in all args
	if ext != "json" && ext != "yaml" {
		// loop through args
		for _, a := range os.Args {
			// get extension
			ext := path.Ext(a)

			// check for valid extensions
			if ext == ".json" || ext == ".yaml" {
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
	if ext == "json" || ext == "yaml" {
		abs := filepath.IsAbs(arg)
		if !abs {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return err
			}

			arg = filepath.Clean(dir + "/" + arg)
		}

		f.FilePath = arg
	} else {
		return errors.New("Error: You must specify a json or yaml file")
	}

	return nil
}

// Load the json/yaml file that was specified from args
// and transform it into a config struct
func (f *File) Load() (*Config, error) {
	var err error
	// read file
	f.Data, err = ioutil.ReadFile(f.FilePath)
	if err != nil {
		return nil, err
	}

	// remove comments
	f.RemoveJSONComments()

	// parse file
	var to *Config
	if f.Mode == "json" {
		err = json.Unmarshal(f.Data, &to)
	} else {
		err = yaml.Unmarshal(f.Data, &to)
	}

	return to, err
}

// RemoveJSONComments from the file
func (f *File) RemoveJSONComments() {
	if f.Mode == "json" {
		// remove inline comments
		f.Data = []byte(regexComments.ReplaceAllString(string(f.Data), ""))
	}
}
