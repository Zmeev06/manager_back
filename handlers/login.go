package handlers

import (
	"stupidauth/models"
	"stupidauth/repos"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx *fiber.Ctx) error {
	var input models.AuthInput
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	user, err := repos.GetUser(input.Login)
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
