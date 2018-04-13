package compactor

import (
	"bytes"
	"compress/gzip"
	"io"
	"time"
)

// gzip compactor
var (
	_ Compactor = NewGZip()
)

type GZip struct{}

func NewGZip() *GZip {
	return new(GZip)
}

func (g *GZip) Name() string {
	return "gzip"
}

func (g *GZip) Encode(src []byte) (dst []byte, err error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	zw.Name = "gzip compactor"
	zw.Comment = "gzip compress"
	zw.ModTime = time.Now()
	if _, err = zw.Write(src); err != nil {
		return
	}

	err = zw.Close()
	dst = buf.Bytes()
	return
}

func (g *GZip) Decode(src []byte) (dst []byte, err error) {
	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(make([]byte, 1024))
	_, err = io.Copy(buf, zr)
	dst = buf.Bytes()
	return
}
