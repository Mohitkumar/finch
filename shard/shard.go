package shard

import (
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/mohitkumar/finch/log"
	"github.com/mohitkumar/finch/storage"
)

type RequestType uint8

const (
	DBPutRequestType     RequestType = 0
	DBDeleteRequestType  RequestType = 1
	LogAppendRequestType RequestType = 2
)

type Shard struct {
	ID      string
	config  Config
	kvStore storage.KVStore
	queues  map[string]log.Log
	raft    *raft.Raft
	dbDir   string
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
	shard.setupStorage()
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
	dbDir := filepath.Join(shardDir, "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}
	shard.dbDir = dbDir
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

func (shard *Shard) setupStorage() {
	config := storage.Config{
		Dir: shard.dbDir,
	}
	shard.kvStore = storage.NewStore(config)
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
