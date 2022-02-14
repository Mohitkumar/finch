package shard

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShard(t *testing.T) {
	f, err := ioutil.TempDir(os.TempDir(), "shard")
	require.NoError(t, err)
	defer os.Remove(f)
}
