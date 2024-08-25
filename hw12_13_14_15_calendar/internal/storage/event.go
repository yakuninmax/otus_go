package storage

import (
	"errors"
	"time"
)

// Event structure.
type Event struct {
	ID          int
	Title       string
	Date        time.Time
	Duration    time.Duration
	Description string
	UserID      int
}

// Error definitions.
var (
	ErrEventNotFound = errors.New("event not found")
)
