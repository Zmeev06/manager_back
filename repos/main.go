package repos

import (
	"os"
	"path"

	"github.com/dgraph-io/badger/v4"
)

var Users *badger.DB
var Keys *badger.DB

var dbs = map[**badger.DB]string{
	&Users: "users", &Keys: "keys"}

const repos = "dbs"

func Init() (err error) {
	if err := os.MkdirAll(repos, 0700); err != nil {
		return err
	}
	for db, name := range dbs {
		db_, err := badger.Open(badger.DefaultOptions(path.Join(repos, name)))
		if err != nil {
			return err
		}
		*db = db_
	}
	return
}
func Close() error {
	for db, _ := range dbs {
		if err := (*db).Close(); err != nil {
			return err
		}
	}
	return nil
}
