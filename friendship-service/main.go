package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/friends/request", func(c *fiber.Ctx) error {
		return c.SendString("Send friend request")
	})

	app.Put("/friends/request/:id", func(c *fiber.Ctx) error {
		return c.SendString("Accept/Reject friend request")
	})

	app.Get("/friends", func(c *fiber.Ctx) error {
		return c.SendString("Get friends list")
	})

	log.Fatal(app.Listen(":3003"))
}
