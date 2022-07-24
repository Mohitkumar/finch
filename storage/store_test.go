package storage

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	f, err := ioutil.TempDir(os.TempDir(), "store_dir")
	require.NoError(t, err)
	defer os.Remove(f)
	c := Config{
		Dir: f,
	}

	store := NewStore(c)

	store.Put([]byte("test"), []byte("value"))
	store.Put([]byte("test1"), []byte("value1"))
	get, err := store.Get([]byte("test"))
	require.NoError(t, err)
	require.Equal(t, "value", string(get))

	file, err := ioutil.TempFile(os.TempDir(), "store_file")
	require.NoError(t, err)
	io.Copy(file, store.Reader())
}
