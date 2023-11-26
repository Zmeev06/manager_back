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
func RmServer(ctx *fiber.Ctx) error {
	var host string
	if err := ctx.BodyParser(&host); err != nil {
		return err
	}
	user := getUserFromJwt(ctx)
	return repos.UpdateUser(user, func(u *models.User) {
		var id int
		for i, v := range u.Servers {
			if v == host {
				id = i
				break
			}
		}
		u.Servers = append(u.Servers[:id], u.Servers[id+1:]...)
	})
}
