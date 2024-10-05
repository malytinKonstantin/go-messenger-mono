package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		// TODO:Здесь будет логика аутентификации и маршрутизации
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Gateway")
	})

	log.Fatal(app.Listen(":3000"))
}
