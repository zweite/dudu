package compactor

type Compactor interface {
	Name() string
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}
