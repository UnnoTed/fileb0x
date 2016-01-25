package template

var fileTemplate = `package {{.Pkg}}

import (
  "log"
  "os"
)

// {{exportedTitle "File"}}{{buildSafeVarName .Path}} is a file
var {{exportedTitle "File"}}{{buildSafeVarName .Path}} = {{.Data}}

func init() {
  f, err := {{exported "FS"}}.OpenFile("{{.Path}}", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  _, err = f.Write({{exportedTitle "File"}}{{buildSafeVarName .Path}})
  if err != nil {
    log.Fatal(err)
  }

  err = f.Close()
  if err != nil {
    log.Fatal(err)
  }
}

`
