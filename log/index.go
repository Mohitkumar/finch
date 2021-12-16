package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

const (
	offsetWidth uint64 = 4
	posWidth    uint64 = 8
	totalWidth  uint64 = offsetWidth + posWidth
)

type (
	index struct {
		file *os.File
		mmap gommap.MMap
		size uint64
	}
)

func newIndex(f *os.File, maxIndexBytes int) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())
	if err := os.Truncate(f.Name(), int64(maxIndexBytes)); err != nil {
		return nil, err
	}

	if idx.mmap, err = gommap.Map(f.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED); err != nil {
		return nil, err
	}
	return idx, nil
}

func (idx *index) Write(offset uint32, pos uint64) error {
	if uint64(len(idx.mmap)) < idx.size+totalWidth {
		return io.EOF
	}
	order.PutUint32(idx.mmap[idx.size:uint64(idx.size+offsetWidth)], offset)
	order.PutUint64(idx.mmap[idx.size+offsetWidth:idx.size+totalWidth], pos)
	idx.size += totalWidth
	return nil
}

func (idx *index) Read(in int64) (pos uint64, off uint32, err error) {
	if idx.size == 0 {
		return 0, 0, io.EOF
	}
	var idxEntry uint32
	if in == -1 {
		idxEntry = uint32((idx.size / totalWidth) - 1)
	} else {
		idxEntry = uint32(in)
	}
	pos = uint64(idxEntry) * totalWidth
	if idx.size < pos+totalWidth {
		return 0, 0, io.EOF
	}

	off = order.Uint32(idx.mmap[pos : pos+offsetWidth])
	pos = order.Uint64(idx.mmap[pos+offsetWidth : pos+totalWidth])
	return pos, off, nil
}

func (idx *index) Close() error {
	if err := idx.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := idx.file.Sync(); err != nil {
		return err
	}
	if err := idx.file.Truncate(int64(idx.size)); err != nil {
		return err
	}
	return idx.file.Close()
}

func (idx *index) Name() string {
	return idx.file.Name()
}
