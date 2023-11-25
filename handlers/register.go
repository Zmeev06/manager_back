package handlers

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	. "stupidauth/models"
)

func Register(ctx *fiber.Ctx) error {
	var net bytes.Buffer
	var input Input
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 5)
	if err != nil {
		return err
	}
	user := User{
		Admin:    false,
		Password: hash,
		Login:    input.Login,
	}
	if user.Login == "admin" {
		user.Admin = true
	}
	if err := gob.NewEncoder(&net).Encode(&user); err != nil {
		return err
	}
	if err := BADGER.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(input.Login)); err == nil {
			return fiber.ErrConflict
		}
		return txn.Set([]byte(input.Login), net.Bytes())
	}); err != nil {
		return err
	}
	str, err := makeToken(user)
	if err != nil {
		return err
	}
	return ctx.JSON(str)
}
