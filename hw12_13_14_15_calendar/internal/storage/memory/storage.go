package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage"
)

// Map for storing events data.
type store map[int]storage.Event

// In memory storage structure.
type Storage struct {
	count  int
	events store
	mu     sync.RWMutex //nolint:unused
}

// Fake storage connection.
func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

// Close fake connection.
func (s *Storage) Close(_ context.Context) {}

// Create event.
func (s *Storage) Create(_ context.Context, event storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	s.events[id] = storage.Event{
		ID:     id,
		Title:  event.Title,
		Date:   event.Date,
		UserID: event.UserID,
	}
	return id, nil
}

// Update event.
func (s *Storage) Update(_ context.Context, id int, change storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get event.
	event, err := s.events[id]
	if !err {
		return storage.ErrEventNotFound
	}

	// Update event data.
	event.Title = change.Title
	event.Date = change.Date
	s.events[id] = event

	return nil
}

// Delete event.
func (s *Storage) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get event.
	_, err := s.events[id]
	if !err {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

// Get events for day.
func (s *Storage) ListDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, day := date.Date()

	// Search events.
	for _, event := range s.events {
		eventYear, eventMonth, eventDay := event.Date.Date()
		if eventYear == year && eventMonth == month && eventDay == day {
			result = append(result, event)
		}
	}

	// Sort by date.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result, nil
}

// Get events for week.
func (s *Storage) ListWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, week := date.ISOWeek()

	// Search events.
	for _, event := range s.events {
		eventYear, eventWeek := event.Date.ISOWeek()
		if eventYear == year && eventWeek == week {
			result = append(result, event)
		}
	}

	// Sort by date.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result, nil
}

// Get events for month.
func (s *Storage) ListMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, _ := date.Date()

	// Search events.
	for _, event := range s.events {
		eventYear, eventMonth, _ := event.Date.Date()
		if eventYear == year && eventMonth == month {
			result = append(result, event)
		}
	}

	// Sort by date.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result, nil
}

// Clean storage.
func (s *Storage) Clean(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(store)
	return nil
}

// Get new event id.
func (s *Storage) newID() int {
	s.count++
	return s.count
}

// Init storage.
func New() *Storage {
	newStorage := Storage{}
	newStorage.events = make(store)
	return &newStorage
}
