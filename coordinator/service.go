package coordinator

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/mohitkumar/finch/discovery"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type CoordinatorService struct {
	Config Config

	mux        cmux.CMux
	coord      *Coordinator
	server     *grpc.Server
	membership *discovery.Membership

	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

func New(config Config) (*CoordinatorService, error) {
	c := &CoordinatorService{
		Config:    config,
		shutdowns: make(chan struct{}),
	}
	setup := []func() error{
		c.setupMux,
		c.setupCoordinator,
		c.setupServer,
		c.setupMembership,
	}
	for _, fn := range setup {
		if err := fn(); err != nil {
			return nil, err
		}
	}
	go c.serve()
	return c, nil
}

func (c *CoordinatorService) setupMux() error {
	rpcAddr := fmt.Sprintf(
		":%d",
		c.Config.RPCPort,
	)
	ln, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}
	c.mux = cmux.New(ln)
	return nil
}

func (c *CoordinatorService) setupCoordinator() error {
	raftLn := c.mux.Match(func(reader io.Reader) bool {
		b := make([]byte, 1)
		if _, err := reader.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(RaftRPC)}) == 0
	})

	c.Config.RaftConfig.StreamLayer = NewStreamLayer(
		raftLn,
	)
	c.Config.RaftConfig.LocalID = raft.ServerID(c.Config.NodeName)
	c.Config.RaftConfig.Bootstrap = c.Config.Bootstrap
	var err error
	c.coord, err = NewCoordinator(c.Config)
	if err != nil {
		return err
	}
	if c.Config.Bootstrap {
		return c.coord.WaitForLeader(3 * time.Second)
	}
	return nil
}

func (c *CoordinatorService) setupMembership() error {
	rpcAddr, err := c.Config.RPCAddr()
	if err != nil {
		return err
	}
	c.membership, err = discovery.New(c.coord, discovery.Config{
		NodeName: c.Config.NodeName,
		BindAddr: c.Config.BindAddr,
		Tags: map[string]string{
			"rpc_addr": rpcAddr,
		},
		StartJoinAddrs: c.Config.StartJoinAddrs,
	})
	return err
}

func (c *CoordinatorService) serve() error {
	if err := c.mux.Serve(); err != nil {
		_ = c.Shutdown()
		return err
	}
	return nil
}

func (c *CoordinatorService) Shutdown() error {
	c.shutdownLock.Lock()
	defer c.shutdownLock.Unlock()
	if c.shutdown {
		return nil
	}
	c.shutdown = true
	close(c.shutdowns)

	shutdown := []func() error{
		c.membership.Leave,
		func() error {
			c.server.GracefulStop()
			return nil
		},
		c.coord.Close,
	}
	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
