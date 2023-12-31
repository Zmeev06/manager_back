package handlers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/gob"

	. "stupidauth/models"
	"stupidauth/repos"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/ssh"
)

func Register(ctx *fiber.Ctx) error {
	// ssh.NewSignerFromKey()
	var net bytes.Buffer
	var input AuthInput
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 5)
	if err != nil {
		return err
	}
	user := User{
		Password: hash,
		Login:    input.Login,
	}
	if user.Login == "admin" {
		user.Admin = true
	}
	if err := gob.NewEncoder(&net).Encode(&user); err != nil {
		return err
	}
	if err := repos.AddUser(user.Login, user); err != nil {
		return err
	}
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	if err := repos.AddKey(user.Login, key); err != nil {
		return err
	}
	str, err := makeToken(user)
	if err != nil {
		return err
	}
	return ctx.JSON(str)
}
