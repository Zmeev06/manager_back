package routing

import (
	"bytes"
	"encoding/gob"
	"stupidauth/handlers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	. "stupidauth/handlers"

	. "stupidauth/models"
	"stupidauth/repos"

	"github.com/dgraph-io/badger/v4"
)

func Setup(app *fiber.App) (err error) {
	err = handlers.Init()
	if err != nil {
		return err
	}
	app.Post("/login", Login)
	app.Post("/reg", Register)
	app.Use(Protected(handlers.JWT_SECRET))
	app.Post("/list_vms", VmList)
	return
}

func Protected(JWT_SECRET []byte) func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: JWT_SECRET},
	})
}
func adminize(ctx *fiber.Ctx) error { // {{{
	var net bytes.Buffer
	login := ctx.Params("user")
	if login == "" {
		return fiber.ErrBadRequest
	}
	return repos.Users.Update(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(login))
		if err != nil {
			return fiber.ErrNotFound
		}
		var user User
		if err := item.Value(func(val []byte) error {
			net.Write(val)
			return gob.NewDecoder(&net).Decode(&user)
		}); err != nil {
			return err
		}
		user.Admin = true
		if err := gob.NewEncoder(&net).Encode(&user); err != nil {
			return err
		}
		return txn.Set([]byte(login), net.Bytes())
	})
} // }}}
