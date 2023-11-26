package handlers

import (
	"stupidauth/models"
	"stupidauth/repos"

	"github.com/gofiber/fiber/v2"
)

func AddServer(ctx *fiber.Ctx) error {
	var host string
	if err := ctx.BodyParser(&host); err != nil {
		return err
	}
	user := getUserFromJwt(ctx)
	return repos.UpdateUser(user, func(u *models.User) {
		u.Servers = append(u.Servers, host)
	})
}
