package template

var fileTemplate = `{{buildTags .Tags}}// Code generaTed by fileb0x at "{{.Now}}" from config file "{{.ConfigFile}}" DO NOT EDIT.
// modified({{.Modified}})
// original path: {{.OriginalPath}}

package {{.Pkg}}

import (
  {{if .Compression.Compress}}
  {{if not .Compression.Keep}}
  "bytes"
  "compress/gzip"
  "io"
  {{end}}
  {{end}}
  "os"
)

// {{exportedTitle "File"}}{{buildSafeVarName .Path}} is "{{.Path}}"
var {{exportedTitle "File"}}{{buildSafeVarName .Path}} = {{.Data}}

func {{ .Init}}() {
  {{if .Compression.Compress}}
  {{if not .Compression.Keep}}
  rb := bytes.NewReader({{exportedTitle "File"}}{{buildSafeVarName .Path}})
  r, err := gzip.NewReader(rb)
  if err != nil {
    panic(err)
  }

  err = r.Close()
  if err != nil {
    panic(err)
  }
  {{end}}
  {{end}}

  f, err := {{exported "FS"}}.OpenFile({{exported "CTX"}}, "{{.Path}}", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    panic(err)
  }

  {{if .Compression.Compress}}
  {{if not .Compression.Keep}}
  _, err = io.Copy(f, r)
  if err != nil {
    panic(err)
  }
  {{end}}
  {{else}}
  _, err = f.Write({{exportedTitle "File"}}{{buildSafeVarName .Path}})
  if err != nil {
    panic(err)
  }
  {{end}}

  err = f.Close()
  if err != nil {
    panic(err)
  }
}

`
