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
	"go.uber.org/zap"
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
	logger       *zap.Logger
}

func New(config Config) (*CoordinatorService, error) {
	c := &CoordinatorService{
		Config:    config,
		shutdowns: make(chan struct{}),
		logger:    zap.L().Named("coordinator-service"),
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
	c.logger.Debug("setting up mux")
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
	c.logger.Debug("setting up coordinator raft")
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
	c.logger.Debug("setting up serf")
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

func (c *CoordinatorService) setupServer() error {
	c.logger.Debug("setting up grpc server for coordinator")
	serverConfig := &GrpcConfig{
		Coordinator: c.coord,
		GetServerer: c.coord,
	}
	var err error
	c.server, err = NewServer(serverConfig)
	if err != nil {
		return err
	}
	grpcLn := c.mux.Match(cmux.Any())
	go func() {
		if err := c.server.Serve(grpcLn); err != nil {
			_ = c.Shutdown()
		}
	}()
	return err
}

func (c *CoordinatorService) serve() error {
	c.logger.Debug("starting coordinator service")
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
