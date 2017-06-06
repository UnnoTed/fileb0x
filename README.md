fileb0x [![Circle CI](https://circleci.com/gh/UnnoTed/fileb0x.svg?style=svg)](https://circleci.com/gh/UnnoTed/fileb0x) [![GoDoc](https://godoc.org/github.com/UnnoTed/fileb0x?status.svg)](https://godoc.org/github.com/UnnoTed/fileb0x) [![GoReportCard](https://goreportcard.com/badge/unnoted/fileb0x)](https://goreportcard.com/report/unnoted/fileb0x)
-------

### What is fileb0x?
A better customizable tool to embed files in go.

It is an alternative to `go-bindata` that have better features and organized configuration.

###### TL;DR
a better `go-bindata`

-------
### How does it compare to `go-bindata`?
Feature                               | fileb0x                       | go-bindata
---------------------                 | -------------                 | ------------------
gofmt                                 | yes (optional)                | no
golint                                | safe                          | unsafe
gzip compression                      | yes                           | yes
gzip decompression                    | yes (optional: runtime)       | yes (on read)
gzip compression levels               | yes                           | no
separated prefix / base for each file | yes                           | no (all files only)
different build tags for each file    | yes                           | no
exclude / ignore files                | yes (glob)                    | yes (regex)
spread files                          | yes                           | no (single file only)
unexported vars/funcs                 | yes (optional)                | no
virtual memory file system            | yes                           | no
http file system / handler            | yes                           | no
replace text in files                 | yes                           | no
glob support                          | yes                           | no (walk folders only)
regex support                         | no                            | yes (ignore files only)
config file                           | yes (config file only)        | no (cmd args only)
update files remotely                 | yes                           | no

-------
### What are the benefits of using a Virtual Memory File System?
By using a virtual memory file system you can have access to files like when they're stored in a hard drive instead of a `map[string][]byte` you would be able to use IO writer and reader.
This means you can `read`, `write`, `remove`, `stat` and `rename` files also `make`, `remove` and `stat` directories.

###### TL;DR 
Virtual Memory File System has similar functions as a hdd stored files would have.



### Features

- [x] golint safe code output

- [x] optional: gzip compression (with optional run-time decompression)

- [x] optional: formatted code (gofmt)

- [x] optional: spread files

- [x] optional: unexporTed variables, functions and types

- [x] optional: include multiple files and folders

- [x] optional: exclude files or/and folders

- [x] optional: replace text in files

- [x] optional: custom base and prefix path

- [x] Virtual Memory FileSystem - [webdav](https://godoc.org/golang.org/x/net/webdav)

- [x] HTTP FileSystem and Handler

- [x] glob support - [doublestar](https://github.com/bmatcuk/doublestar)

- [x] json / yaml / toml support

- [x] optional: Update files remotely

- [x] optional: Build tags for each file


### License
MIT


### Get Started

<details> 

<summary>How to use it?</summary>

##### 1. Download

```bash
go get -u github.com/UnnoTed/fileb0x
```

##### 2. Create a config file
First you need to create a config file, it can be `*.json`, `*.yaml` or `*.toml`. (`*` means any file name)

Now write into the file the configuration you wish, you can use the example files as a start.

json config file example [b0x.json](https://raw.githubusercontent.com/UnnoTed/fileb0x/master/_example/simple/b0x.json)

yaml config file example [b0x.yaml](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/b0x.yaml)

toml config file example [b0x.toml](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/b0x.toml)

##### 3. Run
if you prefer to use it from the `cmd or terminal` edit and run the command below.

```bash
fileb0x YOUR_CONFIG_FILE.yaml
```

or if you wish to generate the embedded files through `go generate` just add and edit the line below into your `main.go`.
```go
//go:generate fileb0x YOUR_CONFIG_FILE.yaml
```

</details>

<details> 
  <summary>What functions and variables fileb0x let me access and what are they for?</summary>

#### HTTP
```go
var HTTP http.FileSystem
```

##### Type 
[`http.FileSystem`](https://golang.org/pkg/net/http/#FileSystem)

##### What is it?

A In-Memory HTTP File System.

##### What it does?

Serve files through a HTTP FileServer.

##### How to use it?
```go
// http.ListenAndServe will create a server at the port 8080
// it will take http.FileServer() as a param
//
// http.FileServer() will use HTTP as a file system so all your files
// can be avialable through the port 8080
http.ListenAndServe(":8080", http.FileServer(myEmbeddedFiles.HTTP))
```
</details>
<details> 
  <summary>How to use it with `echo`?</summary>

```go
package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	// your embedded files import here ...
	"github.com/UnnoTed/fileb0x/_example/echo/myEmbeddedFiles"
)

func main() {
	e := echo.New()

	// enable any filename to be loaded from in-memory file system
	e.GET("/*", echo.WrapHandler(myEmbeddedFiles.Handler))

	// http://localhost:1337/public/README.md
	e.Start(":1337")
}
```

##### How to serve a single file through `echo`?
```go
package main

import (
	"github.com/labstack/echo"

	// your embedded files import here ...
	"github.com/UnnoTed/fileb0x/_example/echo/myEmbeddedFiles"
)

func main() {
	e := echo.New()

	// read ufo.html from in-memory file system
	htmlb, err := myEmbeddedFiles.ReadFile("ufo.html")
	if err != nil {
		log.Fatal(err)
	}

	// convert to string
	html := string(htmlb)

	// serve ufo.html through "/"
	e.GET("/", func(c echo.Context) error {

		// serve as html
		return c.HTML(http.StatusOK, html)
	})

	e.Start(":1337")
}
```

</details>

<details> 
  <summary>Examples</summary>

[simple example](https://github.com/UnnoTed/fileb0x/tree/master/_example/simple) -
[main.go](https://github.com/UnnoTed/fileb0x/blob/master/_example/simple/main.go)

[echo example](https://github.com/UnnoTed/fileb0x/tree/master/_example/echo) -
[main.go](https://github.com/UnnoTed/fileb0x/blob/master/_example/echo/main.go)

```go
package main

import (
	"log"
	"net/http"

  // your generaTed package
	"github.com/UnnoTed/fileb0x/_example/simple/static"
)

func main() {
	files, err := static.WalkDirs("", false)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ALL FILES", files)

  // here we'll read the file from the virtual file system
	b, err := static.ReadFile("public/README.md")
	if err != nil {
		log.Fatal(err)
	}

  // byte to str
  s := string(b)
  s += "#hello"

  // write file back into the virtual file system
  err := static.WriteFile("public/README.md", []byte(s), 0644)
  if err != nil {
    log.Fatal(err)
  }


	log.Println(string(b))

	// true = handler
	// false = file system
	as := false

	// try it -> http://localhost:1337/public/secrets.txt
	if as {
		// as Handler
		panic(http.ListenAndServe(":1337", static.Handler))
	} else {
		// as File System
		panic(http.ListenAndServe(":1337", http.FileServer(static.HTTP)))
	}
}
```
</details>
<details> 

<summary>Update files remotely</summary>

Having to upload an entire binary just to update some files in a b0x and restart a server isn't something that i like to do...

##### How it works?
By enabling the updater option, the next time that you generate a b0x, it will include a http server, this http server will use a http basic auth and it contains 1 endpoint `/` that accepts 2 methods: `GET, POST`.

The `GET` method responds with a list of file names and sha256 hash of each file.
The `POST` method is used to upload files, it creates the directory tree of a new file and then creates the file or it updates an existing file from the virtual memory file system... it responds with a `ok` string when the upload is successful.
  
##### How to update files remotely?

1. First enable the updater option in your config file:
```yaml
##################
## yaml example ##
##################

# updater allows you to update a b0x in a running server
# without having to restart it
updater:
  # disabled by default
  enabled: false

  # empty mode creates a empty b0x file with just the 
  # server and the filesystem, then you'll have to upload
  # the files later using the cmd:
  # fileb0x -update=http://server.com:port b0x.yaml
  #
  # it avoids long compile time
  empty: false

  # amount of uploads at the same time
  workers: 3

  # to get a username and password from a env variable
  # leave username and password blank (username: "")
  # then set your username and password in the env vars 
  # (no caps) -> fileb0x_username and fileb0x_password
  #
  # when using env vars, set it before generating a b0x 
  # so it can be applied to the updater server.
  username: "user" # username: ""
  password: "pass" # password: ""
  port: 8041
```
2. Generate a b0x with the updater option enabled, don't forget to set the username and password for authentication.
3. When your files update, just run `fileb0x -update=http://yourServer.com:8041 b0x.toml` to update the files in the running server.
</details>

<details> 
  <summary>Build Tags</summary>

To use build tags for a b0x package just add the tags to the `tags` property in the main object of your config file
```yaml
# default: main
pkg: static

# destination
dest: "./static/"

# build tags for the main b0x.go file
tags: "!linux"
```

You can also have different build tags for a list of files, you must enable the `spread` property in the main object of your config file, then at the `custom` list, choose the set of files which you want a different build tag 
```yaml
# default: main
pkg: static

# destination
dest: "./static/"

# build tags for the main b0x.go file
tags: "windows darwin"

# [spread] means it will make a file to hold all fileb0x data
# and each file into a separaTed .go file
#
# example:
# theres 2 files in the folder assets, they're: hello.json and world.txt
# when spread is activaTed, fileb0x will make a file: 
# b0x.go or [output]'s data, assets_hello.json.go and assets_world.txt.go
#
#
# type: bool
# default: false
spread: true

# type: array of objects
custom:
  # type: array of strings
  - files: 
    - "start_space_ship.exe"

    # build tags for this set of files
    # it will only work if spread mode is enabled
    tags: "windows"

	# type: array of strings
  - files: 
    - "ufo.dmg"

    # build tags for this set of files
    # it will only work if spread mode is enabled
    tags: "darwin"
```

the config above will make:
```yaml
ab0x.go                         # // +build windows darwin

b0xfile_ufo.exe.go              # // +build windows
b0xfile_start_space_ship.bat.go # // +build darwin
```
</details>

### Functions and Variables

<details> 
  <summary>FS (File System)</summary>

```go
var FS webdav.FileSystem
```

##### Type
[`webdav.FileSystem`](https://godoc.org/golang.org/x/net/webdav#FileSystem)

##### What is it?

In-Memory File System.

##### What it does?

Lets you `read, write, remove, stat and rename` files and `make, remove and stat` directories...

##### How to use it?
```go
func main() {

	// you have the following functions available
	// they all control files/dirs from/to the in-memory file system!
	func Mkdir(name string, perm os.FileMode) error
	func OpenFile(name string, flag int, perm os.FileMode) (File, error)
	func RemoveAll(name string) error
	func Rename(oldName, newName string) error
	func Stat(name string) (os.FileInfo, error)
	// you should remove those lines ^

	// 1. creates a directory
	err := myEmbeddedFiles.FS.Mkdir(myEmbeddedFiles.CTX, "assets", 0777)
	if err != nil {
		log.Fatal(err)
	}

	// 2. creates a file into the directory we created before and opens it
	// with fileb0x you can use ReadFile and WriteFile instead of this complicaTed thing
	f, err := myEmbeddedFiles.FS.OpenFile(myEmbeddedFiles.CTX, "assets/memes.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	data := []byte("I are programmer I make computer beep boop beep beep boop")

	// write the data into the file
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	// close the file
	if err1 := f.Close(); err == nil {
		log.Fatal(err1)
	}

	// 3. rename a file
	// can also move files
	err = myEmbeddedFiles.FS.Rename(myEmbeddedFiles.CTX, "assets/memes.txt", "assets/programmer_memes.txt")
	if err != nil {
		log.Fatal(err)
	}

	// 4. checks if the file we renamed exists
	if _, err = myEmbeddedFiles.FS.Stat(myEmbeddedFiles.CTX, "assets/programmer_memes.txt"); os.IsExist(err) {
		// exists!

		// tries to remove the /assets/ directory
		// from the in-memory file system
		err = myEmbeddedFiles.FS.RemoveAll(myEmbeddedFiles.CTX, "assets")
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. checks if the dir we removed exists
	if _, err = myEmbeddedFiles.FS.Stat(myEmbeddedFiles.CTX, "public/"); os.IsNotExist(err) {
		// doesn't exists!
		log.Println("works!")
	}
}
```
</details>
<details> 
  <summary>Handler</summary>

```go
var Handler *webdav.Handler
```

##### Type
[`webdav.Handler`](https://godoc.org/golang.org/x/net/webdav#Handler)

##### What is it?

A HTTP Handler implementation.

##### What it does?

Serve your embedded files.

##### How to use it?
```go
// ListenAndServer will create a http server at port 8080
// and use Handler as a http handler to serve your embedded files
http.ListenAndServe(":8080", myEmbeddedFiles.Handler)
```
</details>

<details> 
  <summary>ReadFile</summary>

```go
func ReadFile(filename string) ([]byte, error)
```

##### Type
[`ioutil.ReadFile`](https://godoc.org/io/ioutil#ReadFile)

##### What is it?

A Helper function to read your embedded files.

##### What it does?

Reads the specified file from the in-memory file system and return it as a byte slice.

##### How to use it?
```go
// it works the same way that ioutil.ReadFile does.
// but it will read the file from the in-memory file system
// instead of the hard disk!
//
// the file name is passwords.txt
// topSecretFile is a byte slice ([]byte)
topSecretFile, err := myEmbeddedFiles.ReadFile("passwords.txt")
if err != nil {
	log.Fatal(err)
}

log.Println(string(topSecretFile))
```
</details>

<details> 
  <summary>WriteFile</summary>

```go
func WriteFile(filename string, data []byte, perm os.FileMode) error
```

##### Type
[`ioutil.WriteFile`](https://godoc.org/io/ioutil#WriteFile)

##### What is it?

A Helper function to write a file into the in-memory file system.

##### What it does?

Writes the `data` into the specified `filename` in the in-memory file system, meaning you embedded a file!

-- IMPORTANT --
IT WON'T WRITE THE FILE INTO THE .GO GENERATED FILE, IT WILL BE TEMPORARY, WHILE YOUR APP IS RUNNING THE FILE WILL BE AVAILABLE,
AFTER IT SHUTDOWN, IT IS GONE.

##### How to use it?
```go
// it works the same way that ioutil.WriteFile does.
// but it will write the file into the in-memory file system
// instead of the hard disk!
//
// the file name is secret.txt
// data should be a byte slice ([]byte)
// 0644 is a unix file permission

data := []byte("jet fuel can't melt steel beams")
err := myEmbeddedFiles.WriteFile("secret.txt", data, 0644)
if err != nil {
	log.Fatal(err)
}
```
</details>

<details> 
  <summary>WalkDirs</summary>

```go
func WalkDirs(name string, includeDirsInList bool, files ...string) ([]string, error) {
```

##### Type
`[]string`

##### What is it?

A Helper function to walk dirs from the in-memory file system.

##### What it does?

Returns a list of files (with option to include dirs) that are currently in the in-memory file system.

##### How to use it?
```go
includeDirsInTheList := false

// WalkDirs returns a string slice with all file paths
files, err := myEmbeddedFiles.WalkDirs("", includeDirsInTheList)
if err != nil {
	log.Fatal(err)
}

log.Println("List of all my files", files)
```

</details>
