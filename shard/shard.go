package shard

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/raft"
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
	dataDir string
	logDir  string
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
	return shard, nil
}

func (shard *Shard) setupDirectories() error {
	baseDir := shard.config.Dir
	shardDir := filepath.Join(baseDir, shard.ID)
	dataDir := filepath.Join(shardDir, "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	shard.dataDir = dataDir
	logDir := filepath.Join(shardDir, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}
	shard.logDir = logDir
	return nil
}

func (shard *Shard) setupStorage() {
	config := storage.Config{
		Dir: shard.dataDir,
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
