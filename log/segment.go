package log

import (
	"fmt"
	"os"
	"path"

	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/protobuf/proto"
)

type (
	segment struct {
		store      *store
		index      *index
		config     Config
		baseOffset uint64
		nextOffset uint64
	}
)

func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}

	storeFile, err := os.OpenFile(path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}
	s.store, err = newStore(storeFile)
	if err != nil {
		return nil, err
	}

	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}

	if s.index, err = newIndex(indexFile, int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}
	if _, offset, err := s.index.Read(-1); err != nil {
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(offset) + 1
	}
	return s, nil
}

func (s *segment) Append(record *api.LogRecord) (offset uint64, err error) {
	currentOffset := s.nextOffset
	record.Offset = currentOffset

	msg, err := proto.Marshal(record)
	if err != nil {
		return 0, err
	}
	_, pos, err := s.store.Write(msg)
	if err != nil {
		return 0, err
	}
	if err = s.index.Write(uint32(currentOffset-s.baseOffset), pos); err != nil {
		return 0, err
	}
	s.nextOffset++
	return currentOffset, nil
}

func (s *segment) Read(offset uint64) (*api.LogRecord, error) {
	pos, _, err := s.index.Read(int64(offset - s.baseOffset))
	if err != nil {
		return nil, err
	}

	data, err := s.store.Read(pos)
	if err != nil {
		return nil, err
	}

	record := &api.LogRecord{}
	proto.Unmarshal(data, record)
	return record, nil
}

func (s *segment) IsFull() bool {
	return s.index.size >= s.config.Segment.MaxIndexBytes || s.store.size >= s.config.Segment.MaxStoreBytes
}

func (s *segment) Close() error {
	if err := s.index.Close(); err != nil {
		return err
	}
	if err := s.store.Close(); err != nil {
		return err
	}
	return nil
}

func (s *segment) Delete() error {
	if err := s.Close(); err != nil {
		return err
	}
	if err := os.Remove(s.index.Name()); err != nil {
		return err
	}
	if err := os.Remove(s.store.Name()); err != nil {
		return err
	}
	return nil
}
