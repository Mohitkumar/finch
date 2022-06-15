package storage

import badger "github.com/dgraph-io/badger/v3"

func KeyWithCF(cf string, key []byte) []byte {
	return append([]byte(cf+"_"), key...)
}

func GetCF(db *badger.DB, cf string, key []byte) (val []byte, err error) {
	err = db.View(func(txn *badger.Txn) error {
		val, err = GetCFFromTxn(txn, cf, key)
		return err
	})
	return
}

func GetCFFromTxn(txn *badger.Txn, cf string, key []byte) (val []byte, err error) {
	item, err := txn.Get(KeyWithCF(cf, key))
	if err != nil {
		return nil, err
	}
	val, err = item.ValueCopy(val)
	return
}

func PutCF(engine *badger.DB, cf string, key []byte, val []byte) error {
	return engine.Update(func(txn *badger.Txn) error {
		return txn.Set(KeyWithCF(cf, key), val)
	})
}
