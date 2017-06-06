package custom

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/UnnoTed/fileb0x/compression"
	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/UnnoTed/fileb0x/updater"
	"github.com/UnnoTed/fileb0x/utils"
	"github.com/bmatcuk/doublestar"
)

// SharedConfig holds needed data from config package
// without causing import cycle
type SharedConfig struct {
	Output      string
	Compression *compression.Gzip
	Updater     updater.Config
}

// Custom is a set of files with dedicaTed customization
type Custom struct {
	Files  []string
	Base   string
	Prefix string
	Tags   string

	Exclude []string
	Replace []Replacer
}

var xx = []byte(`\x`)

// Parse the files transforming them into a byte string and inserting the file
// into a map of files
func (c *Custom) Parse(files *map[string]*file.File, dirs **dir.Dir, config *SharedConfig) error {
	to := *files
	dirList := *dirs

	var newList []string
	for _, customFile := range c.Files {
		// get files from glob
		list, err := doublestar.Glob(customFile)
		if err != nil {
			return err
		}

		// insert files from glob into the new list
		newList = append(newList, list...)
	}

	// copy new list
	c.Files = newList

	// 0 files in the list
	if len(c.Files) == 0 {
		return errors.New("No files found")
	}

	// loop through files from glob
	for _, customFile := range c.Files {
		// gives error when file doesn't exist
		if !utils.Exists(customFile) {
			return fmt.Errorf("File [%s] doesn't exist", customFile)
		}

		customFile = utils.FixPath(customFile)
		walkErr := filepath.Walk(customFile, func(fpath string, info os.FileInfo, err error) error {
			if config.Updater.Empty && !config.Updater.IsUpdating {
				log.Println("empty mode")
				return nil
			}

			if err != nil {
				return err
			}

			// only files will be processed
			if info.IsDir() {
				return nil
			}

			originalPath := fpath
			fpath = utils.FixPath(fpath)

			var fixedPath string
			if c.Prefix != "" || c.Base != "" {
				c.Base = strings.TrimPrefix(c.Base, "./")

				if strings.HasPrefix(fpath, c.Base) {
					fixedPath = c.Prefix + fpath[len(c.Base):]
				} else {
					if c.Base != "" {
						fixedPath = c.Prefix + fpath
					}
				}

				fixedPath = utils.FixPath(fixedPath)
			} else {
				fixedPath = utils.FixPath(fpath)
			}

			// check for excluded files
			for _, excludedFile := range c.Exclude {
				m, err := doublestar.Match(c.Prefix+excludedFile, fixedPath)
				if err != nil {
					return err
				}

				if m {
					return nil
				}
			}

			// FIXME
			// prevent including itself (destination file or dir)
			if info.Name() == config.Output {
				return nil
			}
			/*if info.Name() == cfg.Output { ||
			  info.Name() == path.Base(path.Dir(jsonFile)) {
			  return nil
			}*/

			// get file's content
			content, err := ioutil.ReadFile(fpath)
			if err != nil {
				return err
			}

			replaced := false

			// loop through replace list
			for _, r := range c.Replace {
				// check if path matches the pattern from property: file
				matched, err := doublestar.Match(c.Prefix+r.File, fixedPath)
				if err != nil {
					return err
				}

				if matched {
					for pattern, word := range r.Replace {
						content = []byte(strings.Replace(string(content), pattern, word, -1))
						replaced = true
					}
				}
			}

			var (
				buf bytes.Buffer
				f   = file.NewFile()
			)

			// it's way faster to use a buffer as string than use string
			buf.WriteString(`[]byte("`)

			// compress the content
			if config.Compression.Options != nil {
				content, err = config.Compression.Compress(content)
				if err != nil {
					return err
				}
			}

			// it's way faster to loop and slice a string than a byte array
			h := hex.EncodeToString(content)

			// loop through hex string, at each 2 chars
			// it's added into a byte array -> []byte("\x09\x11...")
			for i := 0; i < len(h); i += 2 {
				buf.Write(xx)
				buf.WriteString(h[i : i+2])
			}

			f.OriginalPath = originalPath
			f.ReplacedText = replaced
			f.Data = buf.String() + `")`
			f.Name = info.Name()
			f.Path = fixedPath
			f.Tags = c.Tags

			if _, ok := to[fixedPath]; ok {
				f.Tags = to[fixedPath].Tags
			}

			// insert dir to dirlist so it can be created on b0x's init()
			dirList.Insert(path.Dir(fixedPath))

			// insert file into file list
			to[fixedPath] = f
			return nil
		})

		if walkErr != nil {
			return walkErr
		}
	}

	return nil
}
