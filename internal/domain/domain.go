package domain

import "time"

type Monitor struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	IntervalSec int       `json:"interval"`
	CreatedAt   time.Time `json:"omitempty"`
}
