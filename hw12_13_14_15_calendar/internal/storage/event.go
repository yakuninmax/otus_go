package storage

import (
	"errors"
	"time"
)

// Event structure.
type Event struct {
	ID          int
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Description string
	UserID      int
}

// Error definitions.
var (
	ErrEventNotFound     = errors.New("event not found")
	ErrNoScheduledEvents = errors.New("no scheduled events")
	ErrDateIsBusy        = errors.New("date is busy")
)
