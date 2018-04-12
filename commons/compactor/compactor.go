package compactor

type Compactor interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}
