package main

import (
	"bytes"
	"encoding/json"
	"os"
	"stupidauth/handlers"
	"stupidauth/repos"

	"github.com/dgraph-io/badger/v4"
)

func getUserFromJwt(ctx *fiber.Ctx) (user models.User, err error) {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user, err = getUserByName(claims["identity"].(string))
	if err != nil {
		return
	}
	return
}
func main() {
	f()
}

