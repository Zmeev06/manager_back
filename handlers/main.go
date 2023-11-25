package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/golang-jwt/jwt/v5"
	. "stupidauth/models"
)

var JWT_SECRET []byte
var BADGER *badger.DB

type Input struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Init() (err error) {
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	BADGER, err = badger.Open(badger.DefaultOptions("badger"))
	return
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
