package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return c.SendString("Get user profile")
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return c.SendString("Update user profile")
	})

	app.Get("/users/search", func(c *fiber.Ctx) error {
		return c.SendString("Search users")
	})

	log.Fatal(app.Listen(":3002"))
}
