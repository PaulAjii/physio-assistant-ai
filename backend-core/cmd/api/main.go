package main

import (
    "github.com/gofiber/fiber/v3"
	// "log"
)

func main() {
	app := fiber.New()

	app.Get("/", func (c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/upload-voice", func (c fiber.Ctx) error {
		// TODO: Save File to Disc
		// TODO: Send Signal to AI
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Voice recording received, processing has begun",
		})
	})

	app.Listen(":8080")
}