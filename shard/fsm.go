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
	for name, queue := range f.queues {
		readers[i] = queue.Reader(name)
		i++
	}
	return &fsmSnapshot{
		dbReader:  f.kvStore.Reader(),
		logReader: io.MultiReader(readers...),
	}, nil
}

func (f *fsm) Restore(r io.ReadCloser) error {
	defer r.Close()
	var readBuf []byte
	var err error
	if readBuf, err = ioutil.ReadAll(r); err != nil {
		return err
	}
	var i uint64 = 0
	for {
		snapshotItem := &api.SnapShotItem{}
		if err = proto.Unmarshal(readBuf, snapshotItem); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}

		switch item := snapshotItem.Item.(type) {
		case *api.SnapShotItem_KvItem:
			f.kvStore.Put(item.KvItem.Key, item.KvItem.Value)
		case *api.SnapShotItem_LogItem_:
			queuName := item.LogItem.QueueName
			record := &api.LogRecord{}
			if err = proto.Unmarshal(item.LogItem.LogRecord, record); err != nil {
				return err
			}
			log := f.queues[queuName]
			if i == 0 {
				log.GetConfig().Segment.InitialOffset = record.Offset
				if err := log.Reset(); err != nil {
					return err
				}
			}
			log.Append(record)
		}
		i++
	}

	return nil
}
