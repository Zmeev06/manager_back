package repos

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
)

func AddKey(login string, key *rsa.PrivateKey) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(key); err != nil {
		return err
	}
	return Keys.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(login), buf.Bytes())
	})
}
func GetKey(login string) (key rsa.PrivateKey, err error) {
	var net bytes.Buffer
	return key, Keys.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(login))
		if err != nil {
			return fiber.ErrNotFound
		}
		return item.Value(func(val []byte) error {
			_, err := net.Write(val)
			if err != nil {
				return err
			}
			return gob.NewDecoder(&net).Decode(&key)
		})
	})
}
