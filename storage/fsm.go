package storage

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/hashicorp/raft"
	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/protobuf/proto"
)

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	kvstore KVStore
	logger  *log.Logger
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
	f.logger.Printf("put key=%s,value=%s\n", string(key), string(value))
	return f.kvstore.Put(key, value)
}

func (f *fsm) applyDelete(key []byte) error {
	f.logger.Printf("delete key %s \n", string(key))
	return f.kvstore.Delete(key)
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{
		reader: f.kvstore.Reader(),
	}, nil
}

func (f *fsm) Restore(r io.ReadCloser) error {
	f.logger.Printf("Restore snapshot from FSMSnapshot")
	defer r.Close()

	var (
		readBuf  []byte
		err      error
		keyCount int = 0
	)
	// decode message from protobuf
	f.logger.Printf("Read all data")
	if readBuf, err = ioutil.ReadAll(r); err != nil {
		// read done completely
		f.logger.Printf("Snapshot restore failed")
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
			f.logger.Printf("DecodeMessage failed %v", err)
			return err
		}
		// apply item to store
		f.logger.Printf("Set key %v to %v count: %d", item.Key, item.Value, keyCount)
		err = f.kvstore.Put(item.Key, item.Value)
		if err != nil {
			f.logger.Printf("Snapshot load failed %v", err)
			return err
		}
		keyCount = keyCount + 1
	}

	f.logger.Printf("Restore total %d keys", keyCount)

	return nil
}
