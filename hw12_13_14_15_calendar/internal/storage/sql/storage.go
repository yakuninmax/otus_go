package sqlstorage

import (
	"database/sql"
	"fmt"
	"time"

	// Load postgresql driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { // TODO
	db *sql.DB
}

// Init Storage.
func New() *Storage {
	return &Storage{}
}

// Create database connection.
func (s *Storage) Connect(connection *config.StorageConnection) error {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connection.Host,
		connection.Port,
		connection.User,
		connection.Password,
		connection.Database)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	s.db = db

	return nil
}

// Close database connection.
func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("error closing db connection")
	}

	return nil
}

// Create event.
func (s *Storage) Create(event storage.Event) error {
	query := `INSERT INTO events (title, start_date, end_date, description, user_id)
				VALUES($1, $2, $3, $4, $5) RETURNING id`

	args := []interface{}{event.Title, event.StartDate, event.EndDate, event.Description, event.UserID}

	err := s.db.QueryRow(query, args...)
	if err != nil {
		return fmt.Errorf("event creation failed: %w", err.Err())
	}

	return nil
}

// Update event.
func (s *Storage) Update(id int, change storage.Event) error {
	query := `UPDATE events SET title = $1, start_date = $2, end_date = $3,	description = $4 WHERE id = $5`

	args := []interface{}{change.Title, change.StartDate, change.EndDate, change.Description, id}

	result, err := s.db.Exec(query, args...)
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
func (s *Storage) Delete(id int) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := s.db.Exec(query, id)
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
func (s *Storage) ListDay(date time.Time, userID int) ([]storage.Event, error) {
	query := `SELECT * FROM events
				WHERE (extract(year from start_date) = $1 
				AND extract(month from start_date) = $2 
				AND extract(day from start_date) = $3) 
				AND user_id = $4
				ORDER BY start_date`

	year, month, day := date.Date()

	args := []interface{}{year, month, day, userID}

	list, err := s.listQuery(query, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Get events for week.
func (s *Storage) ListWeek(date time.Time, userID int) ([]storage.Event, error) {
	query := `SELECT * FROM events
				WHERE (extract(isoyear from start_date) = $1 AND extract(week from start_date) = $2) 
				AND user_id = $3
				ORDER BY start_date`

	year, week := date.ISOWeek()

	args := []interface{}{year, week, userID}

	list, err := s.listQuery(query, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Get events for month.
func (s *Storage) ListMonth(date time.Time, userID int) ([]storage.Event, error) {
	query := `SELECT * FROM events
				WHERE (extract(year from start_date) = $1 AND extract(month from start_date) = $2) 
				AND user_id = $3
				ORDER BY start_date`

	year, month, _ := date.Date()

	args := []interface{}{year, month, userID}

	list, err := s.listQuery(query, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Clean storage.
func (s *Storage) Clean() error {
	query := `TRUNCATE TABLE events RESTART IDENTITY`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	return nil
}

func (s *Storage) IsBusy(startDate, endDate time.Time, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM events WHERE (start_date = $1 OR end_date = $1 OR start_date = $2 OR end_date = $2
				OR start_date BETWEEN $1 AND $2 OR end_date BETWEEN $1 AND $2) AND user_id = $3`

	args := []interface{}{startDate, endDate, userID}

	var count int
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}

	err = rows.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}

	return count != 0, nil
}

// Listing queries.
func (s *Storage) listQuery(query string, args ...interface{}) ([]storage.Event, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var result []storage.Event
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.StartDate,
			&event.EndDate,
			&event.Description,
			&event.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("database query failed: %w", err)
		}

		result = append(result, event)
	}

	return result, nil
}
