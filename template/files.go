package template

var filesTemplate = `package {{.Pkg}}

import (
  "bytes"
  "io"
  "log"
  "net/http"
  "os"

  "golang.org/x/net/webdav"
)

// FS is a virtual memory file system
var {{exported "FS"}} = webdav.NewMemFS()

// Handler is used to server files through a http handler
var {{exportedTitle "Handler"}} *webdav.Handler

// HTTP is the http file system
var {{exportedTitle "HTTP"}} http.FileSystem = new({{exported "HTTPFS"}})

// HTTPFS implements http.FileSystem
type {{exported "HTTPFS"}} struct {}

{{if not .Spread}}
{{range .Files}}
// {{exportedTitle "File"}}{{buildSafeVarName .Path}} is a file
var {{exportedTitle "File"}}{{buildSafeVarName .Path}} = {{.Data}}
{{end}}
{{end}}

func init() {
  var err error
{{range $index, $dir := .DirList}}
  {{if ne $dir "./"}}
  err = {{exported "FS"}}.Mkdir("{{$dir}}", 0777)
  if err != nil {
    log.Fatal(err)
  }
  {{end}}
{{end}}

{{if not .Spread}}
  var f webdav.File

  {{range .Files}}
  f, err = {{exported "FS"}}.OpenFile("{{.Path}}", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
  {{end}}
{{end}}

  {{exportedTitle "Handler"}} = &webdav.Handler{
    FileSystem: FS,
    LockSystem: webdav.NewMemLS(),
  }
}

// Open a file
func (hfs *{{exported "HTTPFS"}}) Open(path string) (http.File, error) {
  f, err := {{exported "FS"}}.OpenFile(path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  return f, nil
}

// ReadFile is adapTed from ioutil
func {{exportedTitle "ReadFile"}}(path string) ([]byte, error) {
  f, err := {{exported "FS"}}.OpenFile(path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

  // If the buffer overflows, we will get bytes.ErrTooLarge.
  // Return that as an error. Any other panic remains.
  defer func() {
    e := recover()
    if e == nil {
      return
    }
    if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
      err = panicErr
    } else {
      panic(e)
    }
  }()
  _, err = buf.ReadFrom(f)
  return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func {{exportedTitle "WriteFile"}}(filename string, data []byte, perm os.FileMode) error {
  f, err := {{exported "FS"}}.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
  if err != nil {
    return err
  }
  n, err := f.Write(data)
  if err == nil && n < len(data) {
    err = io.ErrShortWrite
  }
  if err1 := f.Close(); err == nil {
    err = err1
  }
  return err
}

`
