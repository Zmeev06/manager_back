package repos

import (
	"bytes"
	"encoding/gob"
	"os"
	"path"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
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
	for db := range dbs {
		if err := (*db).Close(); err != nil {
			return err
		}
	}
	return nil
}
func MakeUpdateFunc[T any](db *badger.DB) func(string, func(*T)) error {
	return func(key string, fn func(*T)) error {
		var net bytes.Buffer
		return db.Update(func(txn *badger.Txn) error {

			resp, err := txn.Get([]byte(key))
			if err != nil {
				return fiber.ErrNotFound
			}
			var item T
			if err := resp.Value(func(val []byte) error {
				net.Write(val)
				return gob.NewDecoder(&net).Decode(&item)
			}); err != nil {
				return err
			}
			fn(&item)
			if err := gob.NewEncoder(&net).Encode(&item); err != nil {
				return err
			}
			return txn.Set([]byte(key), net.Bytes())
		})
	}
}
func MakeGetFunc[T any](db *badger.DB) func(string) (T, error) {
	return func(login string) (item T, err error) {
		var net bytes.Buffer
		return item, Users.View(func(txn *badger.Txn) error {
			resp, err := txn.Get([]byte(login))
			if err != nil {
				return fiber.ErrNotFound
			}
			return resp.Value(func(val []byte) error {
				_, err := net.Write(val)
				if err != nil {
					return err
				}
				return gob.NewDecoder(&net).Decode(&item)
			})
		})
	}
}
func MakeAddFunc[T any](db *badger.DB) func(string, T) error {
	return func(key string, item T) error {
		var net bytes.Buffer
		return db.Update(func(txn *badger.Txn) error {

			if err := gob.NewEncoder(&net).Encode(&item); err != nil {
				return err
			}
			return txn.Set([]byte(key), net.Bytes())
		})
	}
}
