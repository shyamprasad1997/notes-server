package db

import (
	"context"

	"github.com/hashicorp/go-memdb"
)

func (db dbImpl) Txn(ctx context.Context, write bool) MemDbTxn {
	return &memDb{txn: db.db.Txn(write)}
}

type MemDbTxn interface {
	Delete(table string, obj interface{}) error
	Get(table string, index string, args ...interface{}) (memdb.ResultIterator, error)
	First(table string, index string, args ...interface{}) (interface{}, error)
	Insert(table string, obj interface{}) error
	Abort()
	Commit()
}

type memDb struct {
	txn *memdb.Txn
}

func (m *memDb) Insert(table string, obj interface{}) error {
	return m.txn.Insert(table, obj)
}

func (m *memDb) Delete(table string, obj interface{}) error {
	return m.txn.Delete(table, obj)
}

func (m *memDb) Get(table string, index string, args ...interface{}) (memdb.ResultIterator, error) {
	return m.txn.Get(table, index, args)
}

func (m *memDb) First(table string, index string, args ...interface{}) (interface{}, error) {
	return m.txn.First(table, index, args)
}

func (m *memDb) Abort() {
	m.txn.Abort()
}

func (m *memDb) Commit() {
	m.txn.Commit()
}
