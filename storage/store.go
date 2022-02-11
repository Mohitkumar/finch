package storage

import (
	"io"
	"sync"

	"github.com/dgraph-io/badger/v3"
	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/protobuf/proto"
)

type KVStore interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close() error
	Reader() io.ReadCloser
}
type Config struct {
	Dir string
}

type kvStoreImpl struct {
	mu sync.RWMutex
	DB *badger.DB
}

type dbReader struct {
	it *badger.Iterator
}

var _ KVStore = (*kvStoreImpl)(nil)

func NewStore(conf Config) *kvStoreImpl {
	db, err := badger.Open(badger.DefaultOptions(conf.Dir))
	if err != nil {
		panic("failed to open Db in specefied directory")
	}
	s := &kvStoreImpl{
		DB: db,
	}
	return s
}

func (st *kvStoreImpl) Put(key []byte, value []byte) (err error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	err = st.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	return err
}

func (st *kvStoreImpl) Get(key []byte) (valCopy []byte, err error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	err = st.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		valCopy = make([]byte, item.ValueSize())
		_, err = item.ValueCopy(valCopy)
		return err
	})
	if err != nil {
		return nil, err
	}
	return valCopy, nil
}

func (st *kvStoreImpl) Delete(key []byte) error {
	return st.DB.Update(func(txn *badger.Txn) error {
		if err1 := txn.Delete(key); err1 != nil {
			return err1
		}
		return nil
	})
}

func (st *kvStoreImpl) Close() error {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.DB.Close()
}

func (st *kvStoreImpl) Reader() io.ReadCloser {
	st.mu.RLock()
	defer st.mu.RUnlock()
	txn := st.DB.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	it := txn.NewIterator(opts)
	it.Rewind()
	return &dbReader{
		it: it,
	}
}

func (o *dbReader) Read(p []byte) (int, error) {
	if !o.it.Valid() {
		return 0, io.EOF
	}
	item := o.it.Item()
	k := item.Key()
	v := make([]byte, item.ValueSize())
	item.ValueCopy(v)
	kvi := &api.KVItem{
		Key:   k,
		Value: v,
	}
	buff, err := proto.Marshal(kvi)
	if err != nil {
		return 0, err
	}
	n := copy(p, buff)
	o.it.Next()
	return n, nil
}

func (o *dbReader) Close() error {
	o.it.Close()
	return nil
}
