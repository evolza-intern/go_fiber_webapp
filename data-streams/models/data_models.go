package models

import "time"

// DataChunk represents a single piece of incoming data
type DataChunk struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
}

// BatchRequest represents a batch of data chunks
type BatchRequest struct {
	ProcessType string      `json:"process_type"`
	Data        []DataChunk `json:"data"`
}

// StreamResponse represents the response for stream processing
type StreamResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ProcessingJob represents a data processing job
type ProcessingJob struct {
	ID             string             `json:"id"`
	Type           string             `json:"type"`
	Status         string             `json:"status"`
	StartTime      time.Time          `json:"start_time"`
	EndTime        *time.Time         `json:"end_time,omitempty"`
	ProcessedCount int                `json:"processed_count"`
	Results        []ProcessingResult `json:"results,omitempty"`
	Errors         []string           `json:"errors,omitempty"`
}

// ProcessingResult represents the result of processing a data chunk
type ProcessingResult struct {
	ChunkID     string                 `json:"chunk_id"`
	ProcessedAt time.Time              `json:"processed_at"`
	Result      map[string]interface{} `json:"result"`
	Success     bool                   `json:"success"`
	Error       string                 `json:"error,omitempty"`
}

// JobStatus represents the current status of a processing job
type JobStatus struct {
	JobID               string     `json:"job_id"`
	Status              string     `json:"status"`
	Progress            float64    `json:"progress"`
	ProcessedCount      int        `json:"processed_count"`
	StartTime           time.Time  `json:"start_time"`
	EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
}
