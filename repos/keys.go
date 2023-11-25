package repos

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"

	"github.com/dgraph-io/badger/v4"
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
