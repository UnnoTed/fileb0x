package main

import (
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/UnnoTed/fileb0x/config"
	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/UnnoTed/fileb0x/template"
	"github.com/UnnoTed/fileb0x/utils"

	// just to install automatically
	_ "golang.org/x/net/webdav"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	var cfg *config.Config

	// create config and try to get b0x file from args
	j := new(config.JSON)
	err = j.FromArg()

	// required info
	if err != nil {
		log.Fatal(err)
	}

	// load b0x file
	cfg, err = j.Load()
	cfg.Defaults()

	files := make(map[string]*file.File)
	dirs := new(dir.Dir)

	// loop through b0x's [custom] objects
	for _, c := range cfg.Custom {
		err = c.Parse(&files, &dirs)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create files template and exec it
	t := new(template.Template)
	t.Set("files")
	t.Variables = struct {
		Pkg     string
		Files   map[string]*file.File
		Spread  bool
		DirList []string
	}{
		Pkg:     cfg.Pkg,
		Files:   files,
		Spread:  cfg.Spread,
		DirList: dirs.Clean(),
	}
	tmpl, err := t.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// create dest folder when it doesn't exists
	if _, err := os.Stat(cfg.Dest); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.Dest, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	// gofmt
	tmpl, err = format.Source(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// write final execuTed template into the destination file
	err = ioutil.WriteFile(cfg.Dest+cfg.Output, tmpl, 0777)
	if err != nil {
		log.Fatal(err)
	}

	// write spreaded files
	if cfg.Spread {
		a := strings.Split(path.Dir(cfg.Dest), "/")
		dirName := a[len(a)-1:][0]

		for _, f := range files {
			a := strings.Split(path.Dir(f.Path), "/")
			fileDirName := a[len(a)-1:][0]

			if dirName == fileDirName {
				continue
			}

			// transform / to _ and some other chars...
			customName := "b0xfile_" + utils.FixName(f.Path) + ".go"

			// creates file template and exec it
			t := new(template.Template)
			t.Set("file")
			t.Variables = struct {
				Pkg  string
				Path string
				Name string
				Dir  [][]string
				Data string
			}{
				Pkg:  cfg.Pkg,
				Path: f.Path,
				Name: f.Name,
				Dir:  dirs.List,
				Data: f.Data,
			}
			tmpl, err := t.Exec()
			if err != nil {
				log.Fatal(err)
			}

			// gofmt
			tmpl, err = format.Source(tmpl)
			if err != nil {
				log.Fatal(err)
			}

			// write final execuTed template into the destination file
			err = ioutil.WriteFile(cfg.Dest+customName, tmpl, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// success
	log.Println("fileb0x:", cfg.Dest+cfg.Output, "writen!")
}
