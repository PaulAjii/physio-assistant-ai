package main

import (
	"context"

	"github.com/PaulAjii/physio-assistant-ai/internal/config"
	"github.com/PaulAjii/physio-assistant-ai/internal/db"
	"github.com/PaulAjii/physio-assistant-ai/internal/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	cfg := config.Load()

	// Open the database before serving: a bad DSN or unreachable Postgres is a
	// fatal misconfiguration, so fail fast with a clear message rather than
	// starting up and erroring on the first request.
	database, err := db.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer database.Close()

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))
	app.Use(logger.New())

	consultationHandler := handlers.NewConsultationHandler(cfg)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	consultationGroup := app.Group("/api/v1/consultation")
	consultationGroup.Get("/stream/:jobId", consultationHandler.StreamResult)
	consultationGroup.Post("/upload", consultationHandler.UploadAudio)
	consultationGroup.Post("/collate-results", consultationHandler.CollateResults)

	assessmentGroup := app.Group("/api/v1/assessments")
	assessmentGroup.Get("/:id", consultationHandler.GetAssessment)
	assessmentGroup.Put("/:id", consultationHandler.UpdateAssessment)

	// Return from Listen on a clean shutdown so the deferred database.Close
	// runs; only a real server error is fatal.
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
