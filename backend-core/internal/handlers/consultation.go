package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/PaulAjii/physio-assistant-ai/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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

	jobId := uuid.New().String()

	utils.Store.Register(jobId)

	go func() {
		// aiBackendUri := os.Getenv("AI_BACKEND_URI")
		aiBackendUri := "http://localhost:5000/ai/process-audio"
		result, err := utils.ForwardAudioToAI(savedFile, aiBackendUri)
		utils.Store.Complete(jobId, result, err)
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"message":    "Audio file received and processing has begun",
		"statusCode": fiber.StatusOK,
		"jobID":      jobId,
	})
}

func (h *ConsultationHandler) StreamResult(c fiber.Ctx) error {
	jobId := c.Params("jobId")

	job := utils.Store.Get(jobId)
	if job == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "Job not found",
			"statusCode": fiber.ErrNotFound.Code,
		})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	<-utils.Store.Wait(jobId)

	completed := utils.Store.Get(jobId)
	utils.Store.Delete(jobId)

	var event string
	if completed.Error != nil {
		payload, _ := json.Marshal(fiber.Map{
			"status":     "error",
			"message":    completed.Error.Error(),
			"statusCode": fiber.ErrInternalServerError.Code,
		})
		event = fmt.Sprintf("event: error\ndata: %s\n\n", string(payload))
		c.WriteString(event)
		return nil
	}

	var result models.AIResponse
	if err := json.Unmarshal(completed.Result, &result); err != nil {
		payload, _ := json.Marshal(fiber.Map{
			"status":     "error",
			"message":    "Failed to parse AI response",
			"statusCode": fiber.ErrInternalServerError.Code,
		})
		event = fmt.Sprintf("event: error\ndata: %s\n\n", string(payload))
		c.WriteString(event)
		return nil
	}

	payload, _ := json.Marshal(fiber.Map{
		"status":        "success",
		"message":       result.Message,
		"statusCode":    fiber.StatusOK,
		"processedData": result.Data,
	})
	event = fmt.Sprintf("event: result\ndata: %s\n\n", string(payload))

	c.WriteString(event)
	return nil
}
