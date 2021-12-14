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
		dir           string
		activeSegment *segment
		segments      []*segment
	}
)
