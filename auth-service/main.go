package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		return c.SendString("Register endpoint")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login endpoint")
	})

	log.Fatal(app.Listen(":3001"))
}
