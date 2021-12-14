package log

import (
	"os"

	"github.com/tysonmote/gommap"
)

type (
	index struct {
		file *os.File
		mmap *gommap.MMap
		size uint64
	}
)
