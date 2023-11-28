package handlers

import "github.com/gofiber/fiber/v2"

func RegenKeyHand(ctx *fiber.Ctx) error {
	return SetKey(getUserFromJwt(ctx))
}
