package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/notifications", func(c *fiber.Ctx) error {
		return c.SendString("Send notification")
	})

	app.Get("/notifications", func(c *fiber.Ctx) error {
		return c.SendString("Get notifications")
	})

	log.Fatal(app.Listen(":3005"))
}
