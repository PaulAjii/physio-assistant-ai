package handlers

import (
	"github.com/PaulAjii/physio-assistant-ai/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type ConsultationHandler struct{}

func NewConsultationHandler() *ConsultationHandler {
	return &ConsultationHandler{}
}

func (h *ConsultationHandler) UploadAudio(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "Failed to retrieve file. Ensure key is 'file'",
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}

	savedFile, err := utils.SaveUploadedFile(file, "uploads")
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    err.Error(),
			"statusCode": fiber.ErrInternalServerError.Code,
		})
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Consultation audio uploaded successfully",
		"statusCode": fiber.StatusCreated,
		"path":       savedFile,
	})
}
