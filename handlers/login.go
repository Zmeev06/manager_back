package handlers

import (
	"bytes"
	"encoding/gob"

	. "stupidauth/models"
	"stupidauth/repos"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *fiber.Ctx) error {
	var net bytes.Buffer
	var input Input
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	var user User
	err := repos.Users.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(input.Login))
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
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password)); err != nil {
		return fiber.ErrUnauthorized
	}
	str, err := makeToken(user)
	if err != nil {
		return err
	}
	return ctx.JSON(str)
}
