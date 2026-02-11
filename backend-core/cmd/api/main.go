package main

import (
	"github.com/PaulAjii/physio-assistant-ai/internal/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	consultationHandler := handlers.NewConsultationHandler()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")
	api.Post("/upload", consultationHandler.UploadAudio)

	app.Listen(":8080")
}
