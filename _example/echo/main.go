//go:generate fileb0x b0x.yaml
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	// your embedded files import here ...
	"github.com/UnnoTed/fileb0x/_example/echo/myEmbeddedFiles"
	"github.com/UnnoTed/open-golang/open"
)

func main() {
	e := echo.New()
	e.SetDebug(true)

	// enable any filename to be loaded from in-memory file system
	e.GET("/*", standard.WrapHandler(myEmbeddedFiles.Handler))

	// read ufo.html from in-memory file system
	htmlb, err := myEmbeddedFiles.ReadFile("ufo.html")
	if err != nil {
		log.Fatal(err)
	}

	// convert to string
	html := string(htmlb)

	// serve ufo.html through "/"
	e.GET("/", func(c echo.Context) error {

		// serve it
		return c.HTML(http.StatusOK, html)
	})

	// try it -> http://localhost:1337/
	// http://localhost:1337/public/ufo.html
	// http://localhost:1337/public/README.md
	open.Run("http://localhost:1337/")
	e.Run(standard.New(":1337"))

}
