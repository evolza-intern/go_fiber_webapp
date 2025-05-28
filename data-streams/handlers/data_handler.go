package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/evolza-intern/go_fiber_webapp/data-streams/models"
	"github.com/evolza-intern/go_fiber_webapp/data-streams/services"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"strings"
)

// StreamDataHandler processes streaming data chunks
func StreamDataHandler(processor *services.DataProcessor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get processing type from query params
		processType := c.Query("type", "default")

		// Create processing job
		jobID := processor.CreateJob(processType)

		// Set response headers for streaming
		c.Set("Content-Type", "application/json")
		c.Set("Transfer-Encoding", "chunked")

		// Process streaming data
		// Read the entire request body into a buffer first
		bodyBytes := c.Body() // Read the entire body into []byte

		go func() {
			defer processor.CompleteJob(jobID)

			// Create a new bufio.Reader from the in-memory body bytes
			reader := bufio.NewReader(strings.NewReader(string(bodyBytes)))

			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("Job %s: Error reading stream: %v", jobID, err) // Added jobID to log
					break
				}

				// Trim leading/trailing whitespace including newlines and carriage returns
				line = bytes.TrimSpace(line)
				if len(line) == 0 { // Skip empty lines
					continue
				}

				// Parse data chunk
				var dataChunk models.DataChunk
				if err := json.Unmarshal(line, &dataChunk); err != nil {
					log.Printf("Job %s: Error parsing data chunk '%s': %v", jobID, string(line), err) // Added jobID to log
					continue
				}

				// Process chunk asynchronously
				processor.ProcessChunk(jobID, dataChunk)
			}
		}()

		return c.JSON(models.StreamResponse{
			JobID:   jobID,
			Status:  "processing",
			Message: "Stream processing started",
		})
	}
}

// BatchDataHandler processes batch data
func BatchDataHandler(processor *services.DataProcessor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request models.BatchRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		// Create processing job
		jobID := processor.CreateJob(request.ProcessType)

		// Process batch data asynchronously
		go func() {
			defer processor.CompleteJob(jobID)

			for _, chunk := range request.Data {
				processor.ProcessChunk(jobID, chunk)
			}
		}()

		return c.JSON(models.StreamResponse{
			JobID:   jobID,
			Status:  "processing",
			Message: "Batch processing started",
		})
	}
}

// GetProcessingStatus returns job status
func GetProcessingStatus(processor *services.DataProcessor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Query("job_id")
		if jobID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "job_id parameter required",
			})
		}

		status := processor.GetJobStatus(jobID)
		return c.JSON(status)
	}
}

// GetProcessingResults returns job results
func GetProcessingResults(processor *services.DataProcessor) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("id")

		results := processor.GetJobResults(jobID)
		if results == nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Job not found or still processing",
			})
		}

		return c.JSON(results)
	}
}
