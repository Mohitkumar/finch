package coordinator

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/mohitkumar/finch/storage"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type RequestType uint8

const (
	DBPutRequestType    RequestType = 0
	DBDeleteRequestType RequestType = 1
)

type RaftConfig struct {
	StreamLayer *StreamLayer
	raft.Config
	Bootstrap bool
}

type Config struct {
	Dir            string
	NodeName       string
	BindAddr       string
	StartJoinAddrs []string
	RPCPort        int
	Bootstrap      bool
	NumShards      int
	Numreplcia     int
	ShardBaseDir   string
	ShardHost      string
	ShardStartPort int
	RaftConfig     RaftConfig
}

func (c Config) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

type Coordinator struct {
	config  Config
	dbDir   string
	raftDir string
	kvStore storage.KVStore
	raft    *raft.Raft
	logger  *zap.Logger
}

func NewCoordinator(config Config) (*Coordinator, error) {

	coord := &Coordinator{
		config: config,
		logger: zap.L().Named("coordinator"),
	}
	if err := coord.setupDirectories(); err != nil {
		return nil, err
	}
	coord.setupStorage()
	if err := coord.setupRaft(); err != nil {
		return nil, err
	}
	return coord, nil
}

func (c *Coordinator) setupDirectories() error {
	baseDir := c.config.Dir

	dbDir := filepath.Join(baseDir, "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}
	c.dbDir = dbDir

	raftDir := filepath.Join(baseDir, "raft")
	if err := os.MkdirAll(raftDir, 0755); err != nil {
		return err
	}
	c.raftDir = raftDir
	return nil
}

func (c *Coordinator) setupStorage() {
	config := storage.Config{
		Dir: c.dbDir,
	}
	c.kvStore = storage.NewStore(config)
}

func (c *Coordinator) setupRaft() error {
	fsm := &fsm{}
	stableStore, err := raftboltdb.NewBoltStore(
		filepath.Join(c.raftDir, "stable"),
	)
	if err != nil {
		return err
	}

	retain := 1
	snapshotStore, err := raft.NewFileSnapshotStore(
		filepath.Join(c.raftDir, "snap"),
		retain,
		os.Stderr,
	)
	if err != nil {
		return err
	}

	transport := raft.NewNetworkTransport(c.config.RaftConfig.StreamLayer, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	config := raft.DefaultConfig()
	config.LocalID = c.config.RaftConfig.LocalID
	if c.config.RaftConfig.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = c.config.RaftConfig.HeartbeatTimeout
	}
	if c.config.RaftConfig.ElectionTimeout != 0 {
		config.ElectionTimeout = c.config.RaftConfig.ElectionTimeout
	}
	if c.config.RaftConfig.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = c.config.RaftConfig.LeaderLeaseTimeout
	}
	if c.config.RaftConfig.CommitTimeout != 0 {
		config.CommitTimeout = c.config.RaftConfig.CommitTimeout
	}
	c.raft, err = raft.NewRaft(
		config,
		fsm,
		stableStore,
		stableStore,
		snapshotStore,
		transport,
	)
	if err != nil {
		return err
	}
	if c.config.RaftConfig.Bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		c.raft.BootstrapCluster(configuration)
	}
	return nil
}

func (c *Coordinator) Join(id, addr string) error {
	configFuture := c.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}
	serverId := raft.ServerID(id)
	serverAddr := raft.ServerAddress(addr)

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == serverId || srv.Address == serverAddr {
			if srv.ID == serverId && srv.Address == serverAddr {
				return nil
			}
			removeFuture := c.raft.RemoveServer(serverId, 0, 0)
			if err := removeFuture.Error(); err != nil {
				return err
			}
		}
	}
	addFuture := c.raft.AddVoter(serverId, serverAddr, 0, 0)
	if err := addFuture.Error(); err != nil {
		return err
	}
	return nil
}

func (c *Coordinator) Leave(id string) error {
	removeFuture := c.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (c *Coordinator) WaitForLeader(timeout time.Duration) error {
	timeoutc := time.After(timeout)
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-timeoutc:
			return fmt.Errorf("timed out")
		case <-ticker.C:
			if leader := c.raft.Leader(); leader != "" {
				return nil
			}
		}
	}
}

func (c *Coordinator) Close() error {
	shutDownFuture := c.raft.Shutdown()
	if err := shutDownFuture.Error(); err != nil {
		return err
	}
	if err := c.kvStore.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Coordinator) apply(reqType RequestType, req proto.Message) (interface{}, error) {
	var buf bytes.Buffer
	_, err := buf.Write([]byte{byte(reqType)})
	if err != nil {
		return nil, err
	}

	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	if err != nil {
		return nil, err
	}
	applyFuture := c.raft.Apply(buf.Bytes(), 10*time.Second)
	if applyFuture.Error() != nil {
		return nil, applyFuture.Error()
	}
	res := applyFuture.Response()

	if err, ok := res.(error); ok {
		return nil, err
	}
	return res, nil
}
