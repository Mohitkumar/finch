package standalone_storage

import (
	"context"

	"github.com/mohitkumar/finch/storage"
)

type StandAloneStorage struct {
	engines *storage.Engines
}

type Config struct {
	DBPath string
}

func NewStandAloneStorage(conf *Config) *StandAloneStorage {
	kvPath := conf.DBPath + "/kv"
	raftPath := conf.DBPath + "/raft"
	storageCoreDb := storage.CreateDB(kvPath, false)
	raftCoreDb := storage.CreateDB(raftPath, true)
	engines := storage.NewEngines(storageCoreDb, raftCoreDb, kvPath, raftPath)
	return &StandAloneStorage{
		engines: engines,
	}
}

func (s *StandAloneStorage) Reader(ctx context.Context) (storage.StorageReader, error) {
	return NewBadgerReader(s.engines.Kv), nil
}

func (s *StandAloneStorage) Write(ctx context.Context, batch []storage.Modify) error {
	writeBatch := &storage.WriteBatch{}
	for _, m := range batch {
		switch m.Data.(type) {
		case storage.Put:
			writeBatch.SetCF(m.Cf(), m.Key(), m.Value())
		case storage.Delete:
			writeBatch.DeleteCF(m.Cf(), m.Key())
		}
	}
	return s.engines.WriteKV(writeBatch)
}
