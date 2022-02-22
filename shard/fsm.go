package shard

import (
	"io"
	"io/ioutil"

	"github.com/hashicorp/raft"
	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/log"
	"google.golang.org/protobuf/proto"
)

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	queues map[string]log.Log
}

func (l *fsm) Apply(record *raft.Log) interface{} {
	data := record.Data
	reqType := RequestType(data[0])

	switch reqType {
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

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	readers := make([]io.Reader, len(f.queues))
	i := 0
	for name, queue := range f.queues {
		readers[i] = queue.Reader(name)
		i++
	}
	return &fsmSnapshot{
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
		logItem := &api.LogItem{}
		if err = proto.Unmarshal(readBuf, logItem); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}

		queuName := logItem.QueueName
		record := logItem.LogRecord
		log := f.queues[queuName]
		if i == 0 {
			log.GetConfig().Segment.InitialOffset = record.Offset
			if err := log.Reset(); err != nil {
				return err
			}
		}
		log.Append(record)

		i++
	}

	return nil
}
