package main

import (
	"log"
	"stupidauth/repos"
	"stupidauth/routing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.ConfigDefault))
	if err := routing.Setup(app); err != nil {
		log.Fatalln(err)
	}
	defer repos.Close()
	log.Fatalln(app.Listen("127.0.0.1:9090"))
}
