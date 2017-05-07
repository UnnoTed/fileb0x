package updater

import (
	"io"
	"os"
	"testing"

	"time"

	"github.com/UnnoTed/fileb0x/file"
	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	svr := Server{}
	type filef struct {
		filename string
		data     []byte
	}

	go svr.Init()
	time.Sleep(500 * time.Millisecond)

	var files []filef

	files = append(files, filef{
		filename: "whoDid911.txt",
		data:     []byte("notMe"),
	})

	for _, fileData := range files {
		f, err := FS.OpenFile(CTX, fileData.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		assert.NoError(t, err)

		n, err := f.Write(fileData.data)
		if err == nil && n < len(fileData.data) {
			err = io.ErrShortWrite
		}
		if err1 := f.Close(); err == nil {
			err = err1
		}

		assert.NoError(t, err)
		svr.Files = append(svr.Files, fileData.filename)
	}

	up := &Updater{
		Server: "http://localhost:8041",
		Auth: Auth{
			Username: "user",
			Password: "pass",
		},
	}

	assert.Empty(t, up.RemoteHashes)

	err := up.Get()
	assert.NoError(t, err)
	assert.NotEmpty(t, up.RemoteHashes)
	assert.Equal(t, "62e37ff222ec1ca377bb41ffb3fdf08860263e9754f26392c0745765af7397c3", up.RemoteHashes["whoDid911.txt"])

	fileList := map[string]*file.File{
		"whoDid911.txt": {
			Name:  "whoDid911.txt",
			Path:  "whoDid911.txt",
			Bytes: []byte("obsama"),
		},
	}

	updatable, err := up.Updatable(fileList)
	assert.NoError(t, err)
	assert.True(t, updatable)

	err = up.UpdateFiles(fileList)
	assert.NoError(t, err)

	//
	err = up.Get()
	assert.NoError(t, err)
	assert.NotEmpty(t, up.RemoteHashes)
	assert.Equal(t, "98b9640c89068b3639679512d12fd48b4ccb40c1e3fd6b3e0bd4d575cef75cd6", up.RemoteHashes["whoDid911.txt"])

	updatable, err = up.Updatable(fileList)
	assert.NoError(t, err)
	assert.False(t, updatable)
}
