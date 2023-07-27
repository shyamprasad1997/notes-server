package db

import (
	"context"
	"notes-server/models"

	"github.com/hashicorp/go-memdb"
)

type DB interface {
	Txn(ctx context.Context, write bool) *memdb.Txn
}

type dbImpl struct {
	db *memdb.MemDB
}

var dbVar *memdb.MemDB

func NewDB() DB {
	return &dbImpl{db: dbVar}
}

func init() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
					"name": {
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
					"email": {
						Name:    "email",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"password": {
						Name:    "password",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Password"},
					},
				},
			},
			"notes": {
				Name: "notes",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Id"},
					},
					"note": {
						Name:    "note",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Note"},
					},
					"created_by": {
						Name:    "created_by",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "CreatedBy"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	txn := db.Txn(true)
	users := []*models.User{
		{Name: "Admin", Email: "admin@accuknox.com", Password: "admin"},
	}
	for _, user := range users {
		if err := txn.Insert("user", user); err != nil {
			panic(err)
		}
	}
	txn.Commit()
	dbVar = db
}
