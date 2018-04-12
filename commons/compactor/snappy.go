package compactor

import "github.com/golang/snappy"

// snappy compress
var (
	_ Compactor = NewSnappy()
)

type Snappy struct{}

func NewSnappy() *Snappy {
	return new(Snappy)
}

func (s *Snappy) Encode(src []byte) (dst []byte, err error) {
	dst = snappy.Encode(nil, src)
	return
}

func (s *Snappy) Decode(src []byte) (dst []byte, err error) {
	return snappy.Decode(nil, src)
}
