package handlers

import (
	"stupidauth/repos"

	"github.com/gofiber/fiber/v2"
)

func UserInfo(ctx *fiber.Ctx) error {
	username := getUserFromJwt(ctx)
	user, err := repos.GetUser(username)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"login":   user.Login,
		"servers": user.Servers})
}
