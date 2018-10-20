//go:generate fileb0x b0x.yml
package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	// your embedded files import here ...
	"github.com/UnnoTed/fileb0x/_example/gin/myEmbeddedFiles"
)

func main() {
	indexData, err := myEmbeddedFiles.ReadFile("templates/index.html")
	if err != nil {
		panic(err)
	}

	tmpl := template.New("")
	tmpl, err = tmpl.Parse(string(indexData))
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	// Note that thanks to the prefix you can't access templates/index.html here
	r.StaticFS("/css", &myEmbeddedFiles.HTTPFS{Prefix: "css"})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}