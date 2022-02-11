package shard

import "github.com/hashicorp/raft"

type Config struct {
	Dir        string
	RaftConfig RaftConfig
	LogConfig  LogConfig
}
type RaftConfig struct {
	raft.Config
	Bootstrap bool
}

type LogConfig struct {
	MaxStoreBytes uint64
	MaxIndexBytes uint64
	InitialOffset uint64
}
