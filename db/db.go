package db

import (
	"context"

	"github.com/hashicorp/go-memdb"
)

func (db dbImpl) Txn(ctx context.Context, write bool) *memdb.Txn {
	return db.db.Txn(write)
}
