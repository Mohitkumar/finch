package coordinator

import (
	"io"

	"github.com/hashicorp/raft"
)

var _ raft.FSMSnapshot = (*fsmSnapshot)(nil)

type fsmSnapshot struct {
	dbReader io.ReadCloser
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	if _, err := io.Copy(sink, f.dbReader); err != nil {
		_ = sink.Cancel()
		return err
	}
	if err := f.dbReader.Close(); err != nil {
		return err
	}
	return sink.Close()
}

func (f *fsmSnapshot) Release() {

}
