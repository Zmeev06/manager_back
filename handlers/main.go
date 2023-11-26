package handlers

import (
	"os"
	"time"

	. "stupidauth/models"
	"stupidauth/repos"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET []byte

func Init() (err error) {
	if err := repos.Init(); err != nil {
		return err
	}
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	return
}

func makeToken(user User) (str string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user.Login
	claims["exp"] = time.Now().Add(time.Hour * 256).Unix()
	str, err = token.SignedString(JWT_SECRET)
	if err != nil {
		return
	}
	return
}
