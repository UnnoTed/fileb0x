package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	tfile    = "./_example/echo/ufo.html"
	hextable = "0123456789abcdef"
)

var (
	xx      = []byte(`\x`)
	content []byte
)

func Method1(data []byte) string {
	var buf bytes.Buffer
	h := hex.EncodeToString(data)

	// loop through hex string, at each 2 chars
	// it's added into a byte array -> []byte("\x09\x11...")
	for i := 0; i < len(h); i += 2 {
		buf.Write(xx)
		buf.WriteString(h[i : i+2])
	}

	return buf.String()
}

func Method2(data []byte) string {
	dst := make([]byte, len(data)*4)

	for i := 0; i < len(data); i++ {
		dst[i*4] = byte('\\')
		dst[i*4+1] = byte('x')
		dst[i*4+2] = hextable[data[i]>>4]
		dst[i*4+3] = hextable[data[i]&0x0f]
	}

	return string(dst)
}

func BenchmarkInit(b *testing.B) {
	if len(content) > 0 {
		return
	}

	c, err := ioutil.ReadFile(tfile)
	if err != nil {
		log.Fatal(err)
	}

	content = c
}

func BenchmarkOldConvert(b *testing.B) {
	for a := 0; a < b.N; a++ {
		_ = Method1(content)
	}
}

func BenchmarkNewConvert(b *testing.B) {
	for a := 0; a < b.N; a++ {
		_ = Method2(content)
	}
}

func TestMethodsEqual(t *testing.T) {
	if len(content) > 0 {
		return
	}

	c, err := ioutil.ReadFile(tfile)
	if err != nil {
		log.Fatal(err)
	}

	content = c

	assert.Equal(t, Method1(content), Method2(content))
}
