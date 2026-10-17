package domain

import "time"

type Monitor struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	IntervalSec int       `json:"interval"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type CheckResult struct {
	ID             int64
	MonitorID      int64
	StatusCode     int
	ResponseTimeMs int
	ErrorMsg       *string
	CheckedAt      time.Time
}
