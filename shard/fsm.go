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
	case DBPutRequestType:
		item := &api.KVItem{}
		if err := proto.Unmarshal(data[1:], item); err != nil {
			return err
		}
		return l.applyPut(item.Key, item.Value)
	case DBDeleteRequestType:
		return l.applyDelete(data[1:])
	case LogAppendRequestType:
		item := &api.LogItem{}
		if err := proto.Unmarshal(data[1:], item); err != nil {
			return err
		}
		return l.applyAppend(item.QueueName, item.LogRecord)
	}
	return nil
}

func (f *fsm) applyAppend(queueName string, record *api.LogRecord) interface{} {
	offset, err := f.queues[queueName].Append(record)
	if err != nil {
		return err
	}
	return &api.ProduceResponse{Offset: offset}
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
		dataRecord := &api.DataRecord{}
		if err = proto.Unmarshal(readBuf, dataRecord); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}

		switch item := dataRecord.Record.(type) {
		case *api.DataRecord_KvItem:
			f.kvStore.Put(item.KvItem.Key, item.KvItem.Value)
		case *api.DataRecord_LogItem:
			queuName := item.LogItem.QueueName
			record := item.LogItem.LogRecord
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
