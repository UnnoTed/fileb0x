package main

import (
	"flag"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/UnnoTed/fileb0x/compression"
	"github.com/UnnoTed/fileb0x/config"
	"github.com/UnnoTed/fileb0x/custom"
	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/UnnoTed/fileb0x/template"
	"github.com/UnnoTed/fileb0x/updater"
	"github.com/UnnoTed/fileb0x/utils"

	// just to install automatically
	_ "github.com/labstack/echo"
	_ "golang.org/x/net/webdav"
)

var (
	err     error
	cfg     *config.Config
	files   = make(map[string]*file.File)
	dirs    = new(dir.Dir)
	cfgPath string

	fUpdate   string
	startTime = time.Now()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// check for updates
	flag.StringVar(&fUpdate, "update", "", "-update=http(s)://host:port - default port: 8041")
	flag.Parse()
	var (
		update = fUpdate != ""
		up     *updater.Updater
	)

	// create config and try to get b0x file from args
	f := new(config.File)
	err = f.FromArg(true)
	if err != nil {
		log.Fatal(err)
	}

	// load b0x file's config
	cfg, err = f.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.Defaults()
	if err != nil {
		log.Fatal(err)
	}

	cfgPath = f.FilePath

	if err := cfg.Updater.CheckInfo(); err != nil {
		log.Fatal(err)
	}

	cfg.Updater.IsUpdating = update

	// creates a config that can be inserTed into custom
	// without causing a import cycle
	sharedConfig := new(custom.SharedConfig)
	sharedConfig.Output = cfg.Output
	sharedConfig.Updater = cfg.Updater
	sharedConfig.Compression = compression.NewGzip()
	sharedConfig.Compression.Options = cfg.Compression

	// loop through b0x's [custom] objects
	for _, c := range cfg.Custom {
		err = c.Parse(&files, &dirs, sharedConfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	// create files template and exec it
	t := new(template.Template)
	t.Set("files")
	t.Variables = struct {
		ConfigFile  string
		Now         string
		Pkg         string
		Files       map[string]*file.File
		Tags        string
		Spread      bool
		DirList     []string
		Compression *compression.Options
		Debug       bool
		Updater     updater.Config
	}{
		ConfigFile:  filepath.Base(cfgPath),
		Now:         time.Now().String(),
		Pkg:         cfg.Pkg,
		Files:       files,
		Tags:        cfg.Tags,
		Spread:      cfg.Spread,
		DirList:     dirs.Clean(),
		Compression: cfg.Compression,
		Debug:       cfg.Debug,
		Updater:     cfg.Updater,
	}

	tmpl, err := t.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// create dest folder when it doesn't exists
	if !utils.Exists(cfg.Dest) {
		err = os.MkdirAll(cfg.Dest, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	// gofmt
	if cfg.Fmt {
		tmpl, err = format.Source(tmpl)
		if err != nil {
			log.Fatal(err)
		}
	}

	// write final execuTed template into the destination file
	err = ioutil.WriteFile(cfg.Dest+cfg.Output, tmpl, 0777)
	if err != nil {
		log.Fatal(err)
	}

	// write spread files
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
				ConfigFile  string
				Now         string
				Pkg         string
				Path        string
				Name        string
				Dir         [][]string
				Tags        string
				Data        string
				Compression *compression.Options
			}{
				ConfigFile:  filepath.Base(cfgPath),
				Now:         time.Now().String(),
				Pkg:         cfg.Pkg,
				Path:        f.Path,
				Name:        f.Name,
				Dir:         dirs.List,
				Tags:        f.Tags,
				Data:        f.Data,
				Compression: cfg.Compression,
			}
			tmpl, err := t.Exec()
			if err != nil {
				log.Fatal(err)
			}

			// gofmt
			if cfg.Fmt {
				tmpl, err = format.Source(tmpl)
				if err != nil {
					log.Fatal(err)
				}
			}

			// write final execuTed template into the destination file
			err = ioutil.WriteFile(cfg.Dest+customName, tmpl, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// success
	log.Printf("fileb0x: took [%dms] to write [%s] from config file [%s] at [%s]",
		time.Since(startTime).Nanoseconds()/1e6, cfg.Dest+cfg.Output,
		filepath.Base(cfgPath), time.Now().String())

	if update {
		if !cfg.Updater.Enabled {
			log.Fatal("fileb0x: The updater is disabled, enable it in your config file!")
		}

		// includes port when not present
		if !strings.HasSuffix(fUpdate, ":"+strconv.Itoa(cfg.Updater.Port)) {
			fUpdate += ":" + strconv.Itoa(cfg.Updater.Port)
		}

		up = &updater.Updater{
			Server: fUpdate,
			Auth: updater.Auth{
				Username: cfg.Updater.Username,
				Password: cfg.Updater.Password,
			},
			Workers: cfg.Updater.Workers,
		}

		// get file hashes from server
		if err := up.Init(); err != nil {
			panic(err)
		}

		// check if an update is available, then updates...
		if err := up.UpdateFiles(files); err != nil {
			panic(err)
		}
	}
}
