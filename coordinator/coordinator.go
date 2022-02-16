package coordinator

import (
	"fmt"
	"net"
	"sync"

	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/mohitkumar/finch/discovery"
)

type Config struct {
	ConsistentHash struct {
		PartitionCount    uint
		ReplicationFactor uint
	}
	NodeName       string
	BindAddr       string
	StartJoinAddrs []string
	RPCPort        int
	Bootstrap      bool
	NumShards      uint
	Numreplcia     uint
}

func (c Config) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

type Coordinator struct {
	mu          sync.RWMutex
	config      Config
	consistent  *consistent.Consistent
	shardToPeer map[string][]string
	membership  *discovery.Membership
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}
func NewCoordinator(config Config) (*Coordinator, error) {
	cfg := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    int(config.ConsistentHash.PartitionCount),
		ReplicationFactor: int(config.ConsistentHash.ReplicationFactor),
		Load:              1.25,
	}
	c := consistent.New(nil, cfg)
	coord := &Coordinator{
		consistent: c,
		config:     config,
	}
	if err := coord.setupMembership(); err != nil {
		return nil, err
	}
	return coord, nil
}

func (c *Coordinator) setupMembership() error {
	rpcAddr, err := c.config.RPCAddr()
	if err != nil {
		return err
	}
	handler := &handler{}
	c.membership, err = discovery.New(handler, discovery.Config{
		NodeName: c.config.NodeName,
		BindAddr: c.config.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: c.config.StartJoinAddrs,
	})
	return err
}

func (c *Coordinator) setupShards() error {

}
