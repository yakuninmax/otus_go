package storage

import (
	"context"
	"errors"
	"time"
)

// Event structure.
type Event struct {
	ID     int
	Title  string
	Date   time.Time
	UserID int
}

// Events interface.
type Events interface {
	Create(ctx context.Context, event Event) (int, error)
	Update(ctx context.Context, id int, change Event) error
	Delete(ctx context.Context, id int) error
	ListDay(ctx context.Context, date time.Time) ([]Event, error)
	ListWeek(ctx context.Context, date time.Time) ([]Event, error)
	ListMonth(ctx context.Context, date time.Time) ([]Event, error)
}

type Connection interface {
	Connect(ctx context.Context) error
	Close() error
}
type Storage interface {
	Connection
	Events
}

// Error definitions.
var (
	ErrDateBusy      = errors.New("date is busy")
	ErrEventNotFound = errors.New("event not found")
)
