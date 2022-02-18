package coordinator

import (
	"io"
	"io/ioutil"

	"github.com/hashicorp/raft"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/storage"
	"google.golang.org/protobuf/proto"
)

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	kvStore storage.KVStore
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
	return &fsmSnapshot{
		dbReader: f.kvStore.Reader(),
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
		kv := &api.KVItem{}
		if err = proto.Unmarshal(readBuf, kv); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}

		f.kvStore.Put(kv.Key, kv.Value)
		i++
	}

	return nil
}
