fileb0x [![Circle CI](https://circleci.com/gh/UnnoTed/fileb0x.svg?style=svg)](https://circleci.com/gh/UnnoTed/fileb0x) [![Coverage Status](https://coveralls.io/repos/github/UnnoTed/fileb0x/badge.svg?branch=master)](https://coveralls.io/github/UnnoTed/fileb0x?branch=master) [![GoDoc](https://godoc.org/github.com/UnnoTed/fileb0x?status.svg)](https://godoc.org/github.com/UnnoTed/fileb0x) [![GoReportCard](http://goreportcard.com/badge/unnoted/fileb0x)](http://goreportcard.com/report/unnoted/fileb0x)
=======
a simple customizable tool to embed files in go

features:

- [x] golint safe code output

- [x] optional: formatTed code (gofmt)

- [x] optional: spread files

- [x] optional: unexported variables, functions and types

- [x] optional: include multiple files and folders

- [x] optional: exclude files or/and folders

- [x] optional: replace text in files

- [x] optional: custom base and prefix path

- [x] Virtual Memory FileSystem - [webdav](https://godoc.org/golang.org/x/net/webdav)

- [x] HTTP FileSystem and Handler

- [x] glob support - [doublestar](https://github.com/bmatcuk/doublestar)

- [x] json / yaml support

### How to use it?

##### download:

```bash
go get -u github.com/UnnoTed/fileb0x 
```

##### run:

json config file example [b0x.json](https://raw.githubusercontent.com/UnnoTed/fileb0x/master/_example/simple/b0x.json)
```bash
fileb0x b0x.json
```
yaml config file example [b0x.yaml](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/b0x.yaml)
```bash
fileb0x b0x.yaml
```

##### use:

Name                  | Type                                                                            | Description
--------------------- | ------------------------------------------------------------------------------- | ------------------
HTTP                  | var - [http.FileSystem](https://golang.org/pkg/net/http/#FileSystem)            | Serve files through a HTTP FileServer [`http.ListenAndServe(":8080", http.FileServer(static.HTTP))`](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/main.go#L28)
FS                    | var - [webdav.FileSystem](https://godoc.org/golang.org/x/net/webdav#FileSystem) | In-Memory File System, you can `read, write, remove, stat and rename` files, `make, remove and stat` directories...
Handler               | var - [http.Handler](https://golang.org/pkg/net/http/#Handler)                  | Serve file through a HTTP Handler `http.ListenAndServe(":8080", static.Handler)`
ReadFile              | func - [ioutil.ReadFile](https://godoc.org/io/ioutil#ReadFile)                  | Works the same way as [`ioutil.ReadFile`](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/main.go#L11) but the file is read from `FS`
WriteFile             | func - [ioutil.WriteFile](https://godoc.org/io/ioutil#WriteFile)                | Works the same way as `ioutil.WriteFile` but the file is written into `FS`


##### example:

[simple example](https://github.com/UnnoTed/fileb0x/tree/master/_example/simple) -
[main.go](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/main.go)
