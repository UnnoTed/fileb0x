package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// JSON holds config file info
type JSON struct {
	FilePath string
	Data     []byte
}

// FromArg gets the json file from args
func (j *JSON) FromArg() error {
	// (last - 1)
	arg := os.Args[len(os.Args)-1:][0]

	// when json file isn't found on last arg
	// it searches for a "b0x.json" string in all args
	if !strings.HasSuffix(arg, ".json") {
		for _, a := range os.Args {
			if a == "b0x.json" {
				arg = "b0x.json"
			}
		}
	}

	// checks if arg ends with json
	if strings.HasSuffix(arg, ".json") {
		abs := filepath.IsAbs(arg)
		if !abs {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return err
			}

			arg = filepath.Clean(dir + "/" + arg)
		}

		j.FilePath = arg
	} else {
		return errors.New("Error: You must specify the source and destination folder or a json file")
	}

	return nil
}

// Load the json file that was specified from args
// and transform it into a config struct
func (j *JSON) Load() (*Config, error) {
	var err error
	j.Data, err = ioutil.ReadFile(j.FilePath)
	if err != nil {
		return nil, err
	}

	j.RemoveComments()

	var to *Config
	err = json.Unmarshal(j.Data, &to)
	if err != nil {
		return nil, err
	}

	return to, nil
}

// RemoveComments from the json file
func (j *JSON) RemoveComments() {
	// remove inline comments from json
	j.Data = []byte(regexComments.ReplaceAllString(string(j.Data), ""))
}
