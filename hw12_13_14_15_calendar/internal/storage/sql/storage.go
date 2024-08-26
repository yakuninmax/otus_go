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
func (s *Storage) Create(ctx context.Context, event storage.Event) (int, error) {
	query := `
		INSERT INTO events (title, start_date, end_date, description, user_id)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`
	args := []interface{}{event.Title, event.StartDate, event.EndDate, event.Description, event.UserID}

	var id int
	err := s.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("database query failed: %w", err)
	}

	return id, nil
}

// Update event.
func (s *Storage) Update(ctx context.Context, id int, change storage.Event) error {
	query := `
		UPDATE events
		SET title = $1,
			start_date = $2,
			end_date = $3,
			description = $4,
		WHERE id = $5;
	`
	args := []interface{}{change.Title, change.StartDate, change.EndDate, change.Description, id}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if count == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

// Delete event.
func (s *Storage) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM events
		WHERE event_id = $1
	`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if count == 0 {
		return storage.ErrEventNotFound
	}

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
func (s *Storage) Clean(ctx context.Context) error {
	query := `TRUNCATE TABLE event RESTART IDENTITY`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	return nil
}

func New() *Storage {
	return &Storage{}
}
