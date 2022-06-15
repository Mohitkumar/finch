package standalone_storage

import (
	badger "github.com/dgraph-io/badger/v3"
	"github.com/mohitkumar/finch/storage"
)

var _ storage.StorageReader = new(badgerReader)

type badgerReader struct {
	db  *badger.DB
	txn *badger.Txn
	itr *storage.BadgerIterator
}

func NewBadgerReader(db *badger.DB) *badgerReader {
	return &badgerReader{
		db: db,
	}
}

func (r *badgerReader) GetCF(cf string, key []byte) ([]byte, error) {
	val, err := storage.GetCF(r.db, cf, key)
	if err != nil {
		return nil, nil
	}
	return val, nil
}

func (r *badgerReader) IterCF(cf string) storage.DBIterator {
	r.txn = r.db.NewTransaction(false)
	r.itr = storage.NewCFIterator(cf, r.txn)
	return r.itr
}

func (r *badgerReader) Close() {
	r.itr.Close()
	r.txn.Discard()
}
