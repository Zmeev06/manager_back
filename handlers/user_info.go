package handlers

import (
	"stupidauth/repos"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
)

func UserInfo(ctx *fiber.Ctx) error {
	username := getUserFromJwt(ctx)
	user, err := repos.GetUser(username)
	if err != nil {
		return err
	}
	key, err := repos.GetKey(username)
	if err != nil {
		return err
	}
	sshkey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"login":   user.Login,
		"servers": user.Servers,
		"pub_key": string(ssh.MarshalAuthorizedKey(sshkey))})
}
