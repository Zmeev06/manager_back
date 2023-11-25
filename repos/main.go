package repos

import (
	"os"
	"path"

	"github.com/dgraph-io/badger/v4"
)

var Users *badger.DB
var Keys *badger.DB

const repos = "dbs"

func Init() (err error) {
	if err := os.MkdirAll(repos, 0700); err != nil {
		return err
	}
	Users, err = badger.Open(badger.DefaultOptions(path.Join(repos, "users")))
	if err != nil {
		return err
	}
	Keys, err = badger.Open(badger.DefaultOptions(path.Join(repos, "keys")))
	if err != nil {
		return err
	}
	return
}
func Close() error {
	if err := Users.Close(); err != nil {
		return err
	}

	if err := Keys.Close(); err != nil {
		return err
	}
	return nil
}
