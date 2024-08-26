package memorystorage

import (
	"context"
	"fmt"
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
	mu     sync.RWMutex
}

// Fake storage connection.
func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

// Close fake connection.
func (s *Storage) Close(_ context.Context) {}

// Create event.
func (s *Storage) Create(ctx context.Context, event storage.Event) (int, error) {
	// Check if date is busy.
	if s.isBusy(ctx, event.StartDate, event.EndDate, event.UserID) {
		return 0, storage.ErrDateIsBusy
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	s.events[id] = storage.Event{
		ID:          id,
		Title:       event.Title,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Description: event.Description,
		UserID:      event.UserID,
	}
	return id, nil
}

// Update event.
func (s *Storage) Update(ctx context.Context, id int, change storage.Event) error {
	// Check if date is busy.
	if s.isBusy(ctx, change.StartDate, change.EndDate, change.UserID) {
		return storage.ErrDateIsBusy
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Get event.
	event, err := s.events[id]
	if !err {
		return storage.ErrEventNotFound
	}

	// Update event data.
	event.Title = change.Title
	event.StartDate = change.StartDate
	event.EndDate = change.EndDate
	event.Description = change.Description
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
func (s *Storage) ListDay(_ context.Context, date time.Time, userID int) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, day := date.Date()

	// Search events.
	for _, event := range s.events {
		eventYear, eventMonth, eventDay := event.StartDate.Date()
		if eventYear == year && eventMonth == month && eventDay == day && userID == event.UserID {
			result = append(result, event)
		}
	}

	// Check for empty list.
	if len(result) == 0 {
		return nil, storage.ErrNoScheduledEvents
	}

	return sortByDate(result), nil
}

// Get events for week.
func (s *Storage) ListWeek(_ context.Context, date time.Time, userID int) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, week := date.ISOWeek()

	// Search events.
	for _, event := range s.events {
		eventYear, eventWeek := event.StartDate.ISOWeek()
		if eventYear == year && eventWeek == week && userID == event.UserID {
			result = append(result, event)
		}
	}

	// Check for empty list.
	if len(result) == 0 {
		return nil, storage.ErrNoScheduledEvents
	}

	return sortByDate(result), nil
}

// Get events for month.
func (s *Storage) ListMonth(_ context.Context, date time.Time, userID int) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, _ := date.Date()

	// Search events.
	for _, event := range s.events {
		eventYear, eventMonth, _ := event.StartDate.Date()
		if eventYear == year && eventMonth == month && userID == event.UserID {
			result = append(result, event)
		}
	}

	// Check for empty list.
	if len(result) == 0 {
		return nil, storage.ErrNoScheduledEvents
	}

	return sortByDate(result), nil
}

// Clean storage.
func (s *Storage) Clean(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(store)
	s.count = 0
	return nil
}

// Get new event id.
func (s *Storage) newID() int {
	s.count++
	return s.count
}

// Check if time is busy.
func (s *Storage) isBusy(_ context.Context, startDate, endDate time.Time, userID int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	year, month, day := startDate.Date()

	// Search events.
	for _, event := range s.events {
		fmt.Println(event.Title)
		eventYear, eventMonth, eventDay := event.StartDate.Date()
		if eventYear == year && eventMonth == month && eventDay == day && userID == event.UserID {
			if startDate.Equal(event.StartDate) || startDate.Equal(event.EndDate) ||
				endDate.Equal(event.StartDate) || endDate.Equal(event.EndDate) ||
				(startDate.After(event.StartDate) && startDate.Before(event.EndDate)) ||
				(endDate.After(event.StartDate) && endDate.Before(event.EndDate)) {
				return true
			}
		}
	}

	return false
}

func sortByDate(result []storage.Event) []storage.Event {
	// Sort by date.
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartDate.Before(result[j].StartDate)
	})

	return result
}

// Init storage.
func New() *Storage {
	newStorage := Storage{}
	newStorage.events = make(store)
	return &newStorage
}
