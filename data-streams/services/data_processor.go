package services

import (
	"github.com/evolza-intern/go_fiber_webapp/data-streams/models"
	"github.com/google/uuid"
	"log"
	"math"
	"strings"
	"sync"
	"time"
)

// DataProcessor handles concurrent data processing
type DataProcessor struct {
	jobs       map[string]*models.ProcessingJob
	jobsMutex  sync.RWMutex
	workerPool chan struct{}
	shutdown   chan struct{}
	wg         sync.WaitGroup
}

// NewDataProcessor creates a new data processor instance
func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		jobs:       make(map[string]*models.ProcessingJob),
		workerPool: make(chan struct{}, 100), // Limit concurrent workers
		shutdown:   make(chan struct{}),
	}
}

// CreateJob creates a new processing job
func (dp *DataProcessor) CreateJob(processType string) string {
	jobID := uuid.New().String()

	job := &models.ProcessingJob{
		ID:             jobID,
		Type:           processType,
		Status:         "running",
		StartTime:      time.Now(),
		ProcessedCount: 0,
		Results:        make([]models.ProcessingResult, 0),
		Errors:         make([]string, 0),
	}

	dp.jobsMutex.Lock()
	dp.jobs[jobID] = job
	dp.jobsMutex.Unlock()

	log.Printf("Created job %s with type %s", jobID, processType)
	return jobID
}

// ProcessChunk processes a single data chunk
func (dp *DataProcessor) ProcessChunk(jobID string, chunk models.DataChunk) {
	// Acquire worker slot
	dp.workerPool <- struct{}{}
	dp.wg.Add(1)

	go func() {
		defer func() {
			<-dp.workerPool // Release worker slot
			dp.wg.Done()
		}()

		select {
		case <-dp.shutdown:
			return
		default:
			result := dp.processDataChunk(chunk)
			dp.addJobResult(jobID, result)
		}
	}()
}

// processDataChunk performs the actual data processing logic
func (dp *DataProcessor) processDataChunk(chunk models.DataChunk) models.ProcessingResult {
	// Simulate processing time
	time.Sleep(time.Millisecond * 100)

	result := models.ProcessingResult{
		ChunkID:     chunk.ID,
		ProcessedAt: time.Now(),
		Success:     true,
	}

	// Processing logic based on chunk type
	switch chunk.Type {
	case "numeric":
		result.Result = dp.processNumericData(chunk.Data)
	case "text":
		result.Result = dp.processTextData(chunk.Data)
	case "sensor":
		result.Result = dp.processSensorData(chunk.Data)
	default:
		result.Result = dp.processGenericData(chunk.Data)
	}

	return result
}

// processNumericData processes numeric data
func (dp *DataProcessor) processNumericData(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	var sum float64
	var count int

	for _, value := range data {
		if num, ok := value.(float64); ok {
			sum += num
			count++
		}
	}

	if count > 0 {
		result["average"] = sum / float64(count)
		result["sum"] = sum
		result["count"] = count
	}

	result["processed_at"] = time.Now().Unix()
	return result
}

// processTextData processes text data
func (dp *DataProcessor) processTextData(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	var totalLength int
	var wordCount int

	for _, value := range data {
		if text, ok := value.(string); ok {
			totalLength += len(text)
			wordCount += len(strings.Fields(text))
		}
	}

	result["total_length"] = totalLength
	result["word_count"] = wordCount
	result["processed_at"] = time.Now().Unix()
	return result
}

// processSensorData processes sensor data
func (dp *DataProcessor) processSensorData(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if temp, ok := data["temperature"].(float64); ok {
		result["temperature_celsius"] = temp
		result["temperature_fahrenheit"] = (temp * 9 / 5) + 32
		result["temperature_status"] = dp.getTemperatureStatus(temp)
	}

	if humidity, ok := data["humidity"].(float64); ok {
		result["humidity_percent"] = humidity
		result["humidity_status"] = dp.getHumidityStatus(humidity)
	}

	result["processed_at"] = time.Now().Unix()
	return result
}

// processGenericData processes generic data
func (dp *DataProcessor) processGenericData(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["field_count"] = len(data)
	result["processed_at"] = time.Now().Unix()

	// Add field types analysis
	fieldTypes := make(map[string]int)
	for _, value := range data {
		switch value.(type) {
		case string:
			fieldTypes["string"]++
		case float64:
			fieldTypes["number"]++
		case bool:
			fieldTypes["boolean"]++
		default:
			fieldTypes["other"]++
		}
	}
	result["field_types"] = fieldTypes

	return result
}

// Helper functions
func (dp *DataProcessor) getTemperatureStatus(temp float64) string {
	if temp < 0 {
		return "freezing"
	} else if temp < 20 {
		return "cold"
	} else if temp < 30 {
		return "moderate"
	} else {
		return "hot"
	}
}

func (dp *DataProcessor) getHumidityStatus(humidity float64) string {
	if humidity < 30 {
		return "dry"
	} else if humidity < 60 {
		return "comfortable"
	} else {
		return "humid"
	}
}

// addJobResult adds a processing result to a job
func (dp *DataProcessor) addJobResult(jobID string, result models.ProcessingResult) {
	dp.jobsMutex.Lock()
	defer dp.jobsMutex.Unlock()

	if job, exists := dp.jobs[jobID]; exists {
		job.Results = append(job.Results, result)
		job.ProcessedCount++
	}
}

// GetJobStatus returns the current status of a job
func (dp *DataProcessor) GetJobStatus(jobID string) models.JobStatus {
	dp.jobsMutex.RLock()
	defer dp.jobsMutex.RUnlock()

	job, exists := dp.jobs[jobID]
	if !exists {
		return models.JobStatus{
			JobID:  jobID,
			Status: "not_found",
		}
	}

	progress := 0.0
	if job.Status == "completed" {
		progress = 100.0
	} else if job.ProcessedCount > 0 {
		progress = math.Min(float64(job.ProcessedCount)*10, 95.0) // Estimate progress
	}

	return models.JobStatus{
		JobID:          jobID,
		Status:         job.Status,
		Progress:       progress,
		ProcessedCount: job.ProcessedCount,
		StartTime:      job.StartTime,
	}
}

// GetJobResults returns the results of a completed job
func (dp *DataProcessor) GetJobResults(jobID string) *models.ProcessingJob {
	dp.jobsMutex.RLock()
	defer dp.jobsMutex.RUnlock()

	if job, exists := dp.jobs[jobID]; exists && job.Status == "completed" {
		return job
	}
	return nil
}

// CompleteJob marks a job as completed
func (dp *DataProcessor) CompleteJob(jobID string) {
	dp.jobsMutex.Lock()
	defer dp.jobsMutex.Unlock()

	if job, exists := dp.jobs[jobID]; exists {
		job.Status = "completed"
		endTime := time.Now()
		job.EndTime = &endTime
		log.Printf("Job %s completed with %d processed chunks", jobID, job.ProcessedCount)
	}
}

// Shutdown gracefully shuts down the processor
func (dp *DataProcessor) Shutdown() {
	close(dp.shutdown)
	dp.wg.Wait()
	log.Println("Data processor shut down gracefully")
}
