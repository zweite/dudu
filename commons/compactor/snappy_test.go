package compactor

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestSnappy(t *testing.T) {
	snappy := NewSnappy()
	tstr := "some boy like me"

	eb, err := snappy.Encode([]byte(tstr))
	require.Nil(t, err)

	db, err := snappy.Decode(eb)
	require.Nil(t, err)
	require.Equal(t, tstr, string(db))
}
