package shard

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/mohitkumar/finch/log"
)

type RequestType uint8

const (
	LogAppendRequestType RequestType = 0
)

type Shard struct {
	ID      string
	config  Config
	queues  map[string]log.Log
	raft    *raft.Raft
	logDir  string
	raftDir string
}

func NewShard(Id string, confg Config) (*Shard, error) {
	shard := &Shard{
		ID:     Id,
		config: confg,
		queues: make(map[string]log.Log),
	}

	if err := shard.setupDirectories(); err != nil {
		return nil, err
	}
	if err := shard.CreateQueue("system"); err != nil {
		return nil, err
	}

	if err := shard.setupRaft(); err != nil {
		return nil, err
	}
	return shard, nil
}

func (shard *Shard) setupDirectories() error {
	baseDir := shard.config.Dir
	shardDir := filepath.Join(baseDir, shard.ID)
	logDir := filepath.Join(shardDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}
	shard.logDir = logDir

	raftDir := filepath.Join(shardDir, "raft")
	if err := os.MkdirAll(raftDir, 0755); err != nil {
		return err
	}
	shard.raftDir = raftDir
	return nil
}

func (shard *Shard) CreateQueue(name string) error {
	segmentConfig := log.Config{
		Segment: log.Segment{
			InitialOffset: shard.config.LogConfig.InitialOffset,
			MaxIndexBytes: shard.config.LogConfig.MaxIndexBytes,
			MaxStoreBytes: shard.config.LogConfig.MaxStoreBytes,
		},
	}
	l, err := log.NewLog(shard.logDir, segmentConfig)
	if err != nil {
		return err
	}
	shard.queues[name] = l
	return nil
}

func (shard *Shard) setupRaft() error {
	fsm := &fsm{}
	stableStore, err := raftboltdb.NewBoltStore(
		filepath.Join(shard.raftDir, "stable"),
	)
	if err != nil {
		return err
	}

	retain := 1
	snapshotStore, err := raft.NewFileSnapshotStore(
		filepath.Join(shard.raftDir, "snap"),
		retain,
		os.Stderr,
	)
	if err != nil {
		return err
	}
	addr, err := net.ResolveTCPAddr("tcp", shard.config.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(shard.config.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}
	config := raft.DefaultConfig()
	config.LocalID = shard.config.RaftConfig.LocalID
	if shard.config.RaftConfig.HeartbeatTimeout != 0 {
		config.HeartbeatTimeout = shard.config.RaftConfig.HeartbeatTimeout
	}
	if shard.config.RaftConfig.ElectionTimeout != 0 {
		config.ElectionTimeout = shard.config.RaftConfig.ElectionTimeout
	}
	if shard.config.RaftConfig.LeaderLeaseTimeout != 0 {
		config.LeaderLeaseTimeout = shard.config.RaftConfig.LeaderLeaseTimeout
	}
	if shard.config.RaftConfig.CommitTimeout != 0 {
		config.CommitTimeout = shard.config.RaftConfig.CommitTimeout
	}
	shard.raft, err = raft.NewRaft(
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
	if shard.config.RaftConfig.Bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		shard.raft.BootstrapCluster(configuration)
	}
	return nil
}

func (shard *Shard) Join(id, addr string) error {
	configFuture := shard.raft.GetConfiguration()
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
			removeFuture := shard.raft.RemoveServer(serverId, 0, 0)
			if err := removeFuture.Error(); err != nil {
				return err
			}
		}
	}
	addFuture := shard.raft.AddVoter(serverId, serverAddr, 0, 0)
	if err := addFuture.Error(); err != nil {
		return err
	}
	return nil
}

func (shard *Shard) Leave(id string) error {
	removeFuture := shard.raft.RemoveServer(raft.ServerID(id), 0, 0)
	return removeFuture.Error()
}

func (shard *Shard) WaitForLeader(timeout time.Duration) error {
	timeoutc := time.After(timeout)
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-timeoutc:
			return fmt.Errorf("timed out")
		case <-ticker.C:
			if leader := shard.raft.Leader(); leader != "" {
				return nil
			}
		}
	}
}

func (shard *Shard) Close() error {
	shutDownFuture := shard.raft.Shutdown()
	if err := shutDownFuture.Error(); err != nil {
		return err
	}
	for _, queue := range shard.queues {
		queue.Close()
	}
	return nil
}

func (shard *Shard) apply(reqType RequestType, req proto.Message) (interface{}, error) {
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
	applyFuture := shard.raft.Apply(buf.Bytes(), 10*time.Second)
	if applyFuture.Error() != nil {
		return nil, applyFuture.Error()
	}
	res := applyFuture.Response()

	if err, ok := res.(error); ok {
		return nil, err
	}
	return res, nil
}
