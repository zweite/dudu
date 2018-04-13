package proxy

import (
	"bufio"
	"dudu/commons/util"
	"os"
	"path/filepath"
)

// 持久器
type Persistor interface {
	Write([]byte) (int, error)
	Flush() error
	Close() error
}

type FilePersistor struct {
	path string
	fi   *os.File
	buf  *bufio.Writer
}

func NewFilePersistor(filePath string) (*FilePersistor, error) {
	dir := filepath.Dir(filePath)
	if err := util.EnsureDir(dir, 0755); err != nil {
		return nil, err
	}

	fi, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FilePersistor{
		path: filePath,
		fi:   fi,
		buf:  bufio.NewWriter(fi),
	}, nil
}

func (f *FilePersistor) Write(data []byte) (int, error) {
	return f.buf.Write(data)
}

func (f *FilePersistor) Flush() error {
	return f.buf.Flush()
}

func (f *FilePersistor) Close() error {
	f.buf.Flush()
	return f.fi.Close()
}
