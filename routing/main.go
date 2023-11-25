package routing

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET []byte
var BADGER *badger.DB

func Setup(app *fiber.App) (err error) {
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	BADGER, err = badger.Open(badger.DefaultOptions("badger"))
	if err != nil {
		return err
	}
	app.Post("/login", login)
	app.Post("/reg", register)
	app.Post("/adminize/:user", adminize)
	return nil
}

type Input struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func login(ctx *fiber.Ctx) error {
	var net bytes.Buffer
	var input Input
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	var user User
	err := BADGER.View(func(txn *badger.Txn) error {
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
func register(ctx *fiber.Ctx) error {
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
func adminize(ctx *fiber.Ctx) error {
	var net bytes.Buffer
	login := ctx.Params("user")
	if login == "" {
		return fiber.ErrBadRequest
	}
	return BADGER.Update(func(txn *badger.Txn) error {

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
}

type User struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
	Admin    bool   `json:"admin"`
}

func makeToken(user User) (str string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = fmt.Sprint(user.Admin)
	claims["exp"] = time.Now().Add(time.Hour * 256).Unix()
	str, err = token.SignedString(JWT_SECRET)
	if err != nil {
		return
	}
	return
}
