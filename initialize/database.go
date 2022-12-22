package initialize

import (
	"github.com/dgraph-io/badger/v3"
)

func BadgerDB(dirname string) *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(dirname).WithLoggingLevel(badger.ERROR))
	if err != nil {
		panic(err)
	}
	return db
}
