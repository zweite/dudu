package compactor

type Compactor interface {
	Name() string
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

func GetCompactorSet() map[string]Compactor {
	compactorSet := make(map[string]Compactor)
	compactorSet[defaultSnappyCompactor.Name()] = defaultSnappyCompactor
	compactorSet[defaultGZipCompactor.Name()] = defaultGZipCompactor
	return compactorSet
}
