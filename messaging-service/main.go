package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/messages", func(c *fiber.Ctx) error {
		return c.SendString("Send message")
	})

	app.Get("/messages/:chatId", func(c *fiber.Ctx) error {
		return c.SendString("Get chat messages")
	})

	log.Fatal(app.Listen(":3004"))
}
