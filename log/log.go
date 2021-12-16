package log

import (
	"sync"

	api "github.com/mohitkumar/finch/api/v1"
)

type (
	Log interface {
		Append(record *api.LogRecord) (uint64, error)
		Read(offset uint64) (*api.LogRecord, error)
		Close() error
	}

	logImpl struct {
		mu            sync.RWMutex
		Dir           string
		Config        Config
		activeSegment *segment
		segments      []*segment
	}
)

var _ Log = (*logImpl)(nil)

func NewLog(dir string, c Config) (Log, error) {
	if c.Segment.MaxStoreBytes == 0 {
		c.Segment.MaxStoreBytes = 1024
	}
	if c.Segment.MaxIndexBytes == 0 {
		c.Segment.MaxIndexBytes = 1024
	}
	l := &logImpl{
		Dir:    dir,
		Config: c,
	}
	return l, nil
}

func (log *logImpl) Append(record *api.LogRecord) (uint64, error) {

}

func (log *logImpl) Read(offset uint64) (*api.LogRecord, error) {

}

func (log *logImpl) Close() error {

}
