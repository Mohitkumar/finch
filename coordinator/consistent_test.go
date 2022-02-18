package coordinator

import (
	"testing"

	"github.com/buraksezer/consistent"
	"github.com/stretchr/testify/require"
)

type Memb string

func (m Memb) String() string {
	return string(m)
}
func Test(t *testing.T) {
	cfg := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    int(10),
		ReplicationFactor: int(2),
		Load:              1.25,
	}
	c1 := consistent.New(nil, cfg)
	c2 := consistent.New(nil, cfg)

	c1.Add(Memb("test"))
	c1.Add(Memb("test1"))
	c1.Add(Memb("test2"))
	c1.Add(Memb("test3"))
	c1.Add(Memb("test4"))

	c2.Add(Memb("test2"))
	c2.Add(Memb("test1"))
	c2.Add(Memb("test4"))
	c2.Add(Memb("test3"))
	c2.Add(Memb("test"))

	o1 := c1.LocateKey([]byte("my-key"))
	o2 := c2.LocateKey([]byte("my-key"))

	t.Log(o1.String())
	t.Log(o2.String())
	require.Equal(t, o1.String(), o2.String())
	require.Equal(t, c1.LocateKey([]byte("my-key1")).String(), c2.LocateKey([]byte("my-key1")).String())
}
