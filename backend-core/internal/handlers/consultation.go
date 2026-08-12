package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/PaulAjii/physio-assistant-ai/internal/config"
	"github.com/PaulAjii/physio-assistant-ai/internal/models"
	"github.com/PaulAjii/physio-assistant-ai/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ConsultationHandler struct {
	aiBackendURI string
	uploadDir    string
}

func NewConsultationHandler(cfg config.Config) *ConsultationHandler {
	return &ConsultationHandler{
		aiBackendURI: cfg.AIBackendURI,
		uploadDir:    cfg.UploadDir,
	}
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

	savedFile, err := utils.SaveUploadedFile(file, h.uploadDir)
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
		result, err := utils.ForwardAudioToAI(savedFile, h.aiBackendURI)
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

	var event string
	if completed.Error != nil {
		payload, _ := json.Marshal(fiber.Map{
			"status":     "error",
			"message":    completed.Error.Error(),
			"statusCode": fiber.ErrInternalServerError.Code,
		})
		event = fmt.Sprintf("event: error\ndata: %s\n\n", string(payload))
		c.WriteString(event)
		utils.Store.Delete(jobId)
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
		utils.Store.Delete(jobId)
		return nil
	}

	payload, _ := json.Marshal(fiber.Map{
		"status":        "success",
		"message":       result.Message,
		"statusCode":    fiber.StatusOK,
		"processedData": result.Data,
	})
	event = fmt.Sprintf("event: result\ndata: %s\n\n", string(payload))

	utils.Assessments.SaveDraft(jobId, &result)
	utils.Store.Delete(jobId)

	c.WriteString(event)
	return nil
}

func (h *ConsultationHandler) CollateResults(c fiber.Ctx) error {
	var submission *models.AssessmentSubmission
	if err := c.Bind().Body(&submission); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "Invalid request body. Enure it matches the expected format.",
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}

	draft := utils.Assessments.GetDraft(submission.JobID)
	if draft == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "No draft found for the provided job ID",
			"statusCode": fiber.ErrNotFound.Code,
		})
	}

	collated := utils.CollateAssessment(submission, *draft)

	go func() {
		utils.Assessments.SaveAssessment(collated)
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"message":    "Assessment collated and saved successfully",
		"statusCode": fiber.StatusOK,
		"data":       collated,
	})
}

func (h *ConsultationHandler) GetAssessment(c fiber.Ctx) error {
	id := c.Params("id")

	assessment := utils.Assessments.GetAssessment(id)
	if assessment == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "Assessment not found",
			"statusCode": fiber.ErrNotFound.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"message":    "Assessment retrieved successfully",
		"statusCode": fiber.StatusOK,
		"data":       assessment,
	})
}

func (h *ConsultationHandler) UpdateAssessment(c fiber.Ctx) error {
	id := c.Params("id")

	existing := utils.Assessments.GetAssessment(id)
	if existing == nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "Assessment not found",
			"statusCode": fiber.ErrNotFound.Code,
		})
	}

	var updateData models.Assessment
	if err := c.Bind().Body(&updateData); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":     "error",
			"message":    "Invalid request body",
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}

	existing.Assessment = updateData

	go func() {
		utils.Assessments.SaveAssessment(existing)
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"message":    "Assessment updated successfully",
		"statusCode": fiber.StatusOK,
		"data":       existing,
	})
}
