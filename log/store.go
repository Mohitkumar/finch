package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	order = binary.BigEndian
)

const (
	dataLengthWidth = 8
)

type (
	store struct {
		file *os.File
		mu   sync.Mutex
		buf  *bufio.Writer
		size uint64
	}
)

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := fi.Size()
	return &store{
		file: f,
		size: uint64(size),
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *store) Write(data []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	if err = binary.Write(s.buf, order, uint64(len(data))); err != nil {
		return 0, 0, err
	}
	bytesWritten, err := s.buf.Write(data)
	if err != nil {
		return 0, 0, err
	}
	totalBytesWritten := bytesWritten + dataLengthWidth

	s.size += uint64(totalBytesWritten)
	return uint64(totalBytesWritten), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, dataLengthWidth)
	if _, err := s.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	data := make([]byte, order.Uint64(size))
	if _, err := s.ReadAt(data, int64(pos+dataLengthWidth)); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *store) ReadAt(data []byte, offset int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.file.ReadAt(data, offset)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.file.Close()
}

func (s *store) Name() string {
	return s.file.Name()
}
