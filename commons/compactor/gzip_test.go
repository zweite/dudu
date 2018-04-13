package compactor

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestGZip(t *testing.T) {
	gzip := NewGZip()
	tstr := "some boy like me"

	eb, err := gzip.Encode([]byte(tstr))
	require.Nil(t, err)

	db, err := gzip.Decode(eb)
	require.Nil(t, err)
	require.Equal(t, tstr, string(db))
}
