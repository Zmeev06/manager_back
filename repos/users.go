package repos

import (
	"bytes"
	"encoding/gob"
	"stupidauth/models"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
)

func UpdateUser(login string, fn func(*models.User)) error {
	var net bytes.Buffer
	return Users.Update(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(login))
		if err != nil {
			return fiber.ErrNotFound
		}
		var user models.User
		if err := item.Value(func(val []byte) error {
			net.Write(val)
			return gob.NewDecoder(&net).Decode(&user)
		}); err != nil {
			return err
		}
		fn(&user)
		if err := gob.NewEncoder(&net).Encode(&user); err != nil {
			return err
		}
		return txn.Set([]byte(login), net.Bytes())
	})
}
func GetUser(login string) (user models.User, err error) {
	var net bytes.Buffer
	return user, Users.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(login))
		if err != nil {
			return fiber.ErrNotFound
		}
		return item.Value(func(val []byte) error {
			_, err := net.Write(val)
			if err != nil {
				return err
			}
			return gob.NewDecoder(&net).Decode(&user)
		})
	})
}
