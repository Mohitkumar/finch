package shard

import (
	"io"
	"io/ioutil"

	"github.com/hashicorp/raft"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/log"
	"github.com/mohitkumar/finch/storage"
	"google.golang.org/protobuf/proto"
)

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	kvStore storage.KVStore
	queues  map[string]log.Log
}

func (l *fsm) Apply(record *raft.Log) interface{} {
	data := record.Data
	reqType := RequestType(data[0])

	switch reqType {
	case PutRequestType:
		item := &api.KVItem{}
		if err := proto.Unmarshal(data[1:], item); err != nil {
			return err
		}
		return l.applyPut(item.Key, item.Value)
	case DeleteRequestType:
		return l.applyDelete(data[1:])
	}
	return nil
}

func (f *fsm) applyPut(key []byte, value []byte) error {
	return f.kvStore.Put(key, value)
}

func (f *fsm) applyDelete(key []byte) error {
	return f.kvStore.Delete(key)
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	readers := make([]io.Reader, len(f.queues))
	i := 0
	for _, queue := range f.queues {
		readers[i] = queue.Reader()
		i++
	}
	return &fsmSnapshot{
		dbReader:  f.kvStore.Reader(),
		logReader: io.MultiReader(readers...),
	}, nil
}

func (f *fsm) Restore(r io.ReadCloser) error {
	defer r.Close()

	var (
		readBuf  []byte
		err      error
		keyCount int = 0
	)
	// decode message from protobuf
	if readBuf, err = ioutil.ReadAll(r); err != nil {
		// read done completely
		return err
	}

	// decode messages from 1M block file
	// the last message could decode failed with io.ErrUnexpectedEOF
	for {
		item := &api.KVItem{}
		if err = proto.Unmarshal(readBuf, item); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}
		// apply item to store
		err = f.kvStore.Put(item.Key, item.Value)
		if err != nil {
			return err
		}
		keyCount = keyCount + 1
	}

	return nil
}
