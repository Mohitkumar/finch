package storage

import (
	"io"

	"github.com/hashicorp/raft"
)

var _ raft.FSMSnapshot = (*fsmSnapshot)(nil)

type fsmSnapshot struct {
	reader io.ReadCloser
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	if _, err := io.Copy(sink, f.reader); err != nil {
		_ = sink.Cancel()
		return err
	}
	if err := f.reader.Close(); err != nil {
		return err
	}
	return sink.Close()
}

func (f *fsmSnapshot) Release() {

}
