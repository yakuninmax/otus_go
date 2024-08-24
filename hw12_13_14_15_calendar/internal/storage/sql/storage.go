package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	db *sql.DB
}

// Create database connection.
func (s *Storage) Connect(ctx context.Context, connect string) error {
	db, err := sql.Open("pgx", connect)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	s.db = db

	return s.db.PingContext(ctx)
}

// Close database connection.
func (s *Storage) Close(_ context.Context) {
	s.db.Close()
}

// Create event.
func (s *Storage) Create(_ context.Context, event storage.Event) (int, error) {
	return 0, nil
}

// Update event.
func (s *Storage) Update(_ context.Context, id int, change storage.Event) error {
	return nil
}

// Delete event.
func (s *Storage) Delete(_ context.Context, id int) error {
	return nil
}

// Get events for day.
func (s *Storage) ListDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

// Get events for week.
func (s *Storage) ListWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

// Get events for month.
func (s *Storage) ListMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}

// Clean storage.
func (s *Storage) Clean(_ context.Context) error {
	return nil
}

func New() *Storage {
	return &Storage{}
}
