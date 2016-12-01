package template

var filesTemplate = `package {{.Pkg}}
{{$Compression := .Compression}}

import (
  "bytes"
  {{if not .Spread}}{{if $Compression.Compress}}{{if not $Compression.Keep}}"compress/gzip"{{end}}{{end}}{{end}}
  "io"
  "log"
  "net/http"
  "os"

  "golang.org/x/net/webdav"
  "golang.org/x/net/context"
)

var ( 
  // CTX is a context for webdav vfs
  {{exported "CTX"}} = context.Background()

  {{if .Debug}}
  {{exported "FS"}} = webdav.Dir(".")
  {{else}}
  // FS is a virtual memory file system
  {{exported "FS"}} = webdav.NewMemFS()
  {{end}}

  // Handler is used to server files through a http handler
  {{exportedTitle "Handler"}} *webdav.Handler

  // HTTP is the http file system
  {{exportedTitle "HTTP"}} http.FileSystem = new({{exported "HTTPFS"}})
)

// HTTPFS implements http.FileSystem
type {{exported "HTTPFS"}} struct {}

{{if (and (not .Spread) (not .Debug))}}
{{range .Files}}
// {{exportedTitle "File"}}{{buildSafeVarName .Path}} is a file
var {{exportedTitle "File"}}{{buildSafeVarName .Path}} = {{.Data}}
{{end}}
{{end}}

func init() {
  if {{exported "CTX"}}.Err() != nil {
		log.Fatal({{exported "CTX"}}.Err())
	}

{{if not .Debug}}
var err error
{{range $index, $dir := .DirList}}
  {{if ne $dir "./"}}
  err = {{exported "FS"}}.Mkdir({{exported "CTX"}}, "{{$dir}}", 0777)
  if err != nil {
    log.Fatal(err)
  }
  {{end}}
{{end}}
{{end}}

{{if (and (not .Spread) (not .Debug))}}
  var f webdav.File
  {{if $Compression.Compress}}
  {{if not $Compression.Keep}}
  var rb *bytes.Reader
  var r *gzip.Reader
  {{end}}
  {{end}}

  {{range .Files}}
  {{if $Compression.Compress}}
  {{if not $Compression.Keep}}
  rb = bytes.NewReader({{exportedTitle "File"}}{{buildSafeVarName .Path}})
  r, err = gzip.NewReader(rb)
  if err != nil {
    log.Fatal(err)
  }

  err = r.Close()
  if err != nil {
    log.Fatal(err)
  }
  {{end}}
  {{end}}

  f, err = {{exported "FS"}}.OpenFile({{exported "CTX"}}, "{{.Path}}", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
  if err != nil {
    log.Fatal(err)
  }

  {{if $Compression.Compress}}
  {{if not $Compression.Keep}}
  _, err = io.Copy(f, r)
  if err != nil {
    log.Fatal(err)
  }
  {{end}}
  {{else}}
  _, err = f.Write({{exportedTitle "File"}}{{buildSafeVarName .Path}})
  if err != nil {
    log.Fatal(err)
  }
  {{end}}

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
  f, err := {{if .Debug}}os{{else}}{{exported "FS"}}{{end}}.OpenFile({{if not .Debug}}{{exported "CTX"}}, {{end}}path, os.O_RDONLY, 0644)
  if err != nil {
    return nil, err
  }

  return f, nil
}

// ReadFile is adapTed from ioutil
func {{exportedTitle "ReadFile"}}(path string) ([]byte, error) {
  f, err := {{if .Debug}}os{{else}}{{exported "FS"}}{{end}}.OpenFile({{if not .Debug}}{{exported "CTX"}}, {{end}}path, os.O_RDONLY, 0644)
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
  f, err := {{exported "FS"}}.OpenFile({{exported "CTX"}}, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
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

// FileNames is a list of files included in this filebox
var {{exportedTitle "FileNames"}} = []string {
  {{range .Files}}"{{.Path}}",
  {{end}}
}

`
