package log

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	api "github.com/mohitkumar/finch/api/v1"
)

type (
	Log interface {
		Append(record *api.LogRecord) (uint64, error)
		Read(offset uint64) (*api.LogRecord, error)
		Close() error
		Reader() io.Reader
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
	return l, l.setup()
}

func (log *logImpl) setup() error {
	files, err := ioutil.ReadDir(log.Dir)
	if err != nil {
		return err
	}
	var baseOffsets []uint64
	for _, file := range files {
		offsetStr := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		off, _ := strconv.ParseUint(offsetStr, 10, 0)
		baseOffsets = append(baseOffsets, off)
	}
	sort.Slice(baseOffsets, func(i, j int) bool {
		return baseOffsets[i] < baseOffsets[j]
	})
	for i := 0; i < len(baseOffsets); i++ {
		if err := log.newSegment(baseOffsets[i]); err != nil {
			return err
		}
		i++
	}
	if log.segments == nil {
		if err := log.newSegment(log.Config.Segment.InitialOffset); err != nil {
			return err
		}
	}
	return nil
}

func (log *logImpl) newSegment(offset uint64) error {
	s, err := newSegment(log.Dir, offset, log.Config)
	if err != nil {
		return err
	}
	log.segments = append(log.segments, s)
	log.activeSegment = s
	return nil
}

func (log *logImpl) Append(record *api.LogRecord) (uint64, error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	off, err := log.activeSegment.Append(record)
	if err != nil {
		return 0, err
	}
	if log.activeSegment.IsFull() {
		err = log.newSegment(off + 1)
	}
	return off, err
}

func (log *logImpl) Read(offset uint64) (*api.LogRecord, error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	var s *segment
	for _, segment := range log.segments {
		if offset >= segment.baseOffset && offset < segment.nextOffset {
			s = segment
			break
		}
	}
	if s == nil || s.nextOffset <= offset {
		return nil, api.ErrOffsetOutOfRange{Offset: offset}
	}
	// END: after
	return s.Read(offset)
}

func (l *logImpl) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, segment := range l.segments {
		if err := segment.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (l *logImpl) Remove() error {
	if err := l.Close(); err != nil {
		return err
	}
	return os.RemoveAll(l.Dir)
}

func (l *logImpl) Reset() error {
	if err := l.Remove(); err != nil {
		return err
	}
	return l.setup()
}

func (l *logImpl) Reader() io.Reader {
	l.mu.RLock()
	defer l.mu.RUnlock()
	readers := make([]io.Reader, len(l.segments))
	for i, segment := range l.segments {
		readers[i] = &originReader{segment.store, 0}
	}
	return io.MultiReader(readers...)
}

type originReader struct {
	*store
	off int64
}

func (o *originReader) Read(p []byte) (int, error) {
	n, err := o.ReadAt(p, o.off)
	o.off += int64(n)
	return n, err
}
