package models

import "time"

// ActivityLog represents application activity
type ActivityLog struct {
	ID        int       `json:"id"`
	Action    string    `json:"action"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}
