package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"testing"
)

const (
	tfile = "_example/simple/public/index.html"
)

var (
	content []byte
	h       = hex.EncodeToString(content)
)

func BenchmarkInit(b *testing.B) {
	c, err := ioutil.ReadFile(tfile)
	if err != nil {
		log.Fatal(err)
	}

	content = c
}

func BenchmarkConvert(b *testing.B) {
	for a := 0; a < b.N; a++ {
		var buf bytes.Buffer
		buf.WriteString(`[]byte("`)

		h := hex.EncodeToString(content)
		// loop through hex string, at each 2 chars
		// it's added into a byte array -> []byte{0x61 ,...}
		for i := 0; i < len(h); i += 2 {
			buf.WriteString("\\x" + h[i:i+2])
		}

		// remove last comma and insert -> }
		_ = buf.String() + `")`
	}
}

func TestCreation(t *testing.T) {

}
