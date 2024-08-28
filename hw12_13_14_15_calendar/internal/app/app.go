package app

import (
	"errors"
	"time"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage"
)

// App structure.
type App struct {
	Logger
	Storage
}

// Logger interface.
type Logger interface {
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debug(...interface{})
}

var (
	ErrNoUserID       = errors.New("user id not set")
	ErrNoEventID      = errors.New("event id not set")
	ErrEmptyTitle     = errors.New("empty title")
	ErrNoStartDate    = errors.New("start date is empty")
	ErrNoEndDate      = errors.New("end date is empty")
	ErrEndBeforeStart = errors.New("end date earlier than start date")
	ErrTimeBusy       = errors.New("time is busy")
)

// Storage interface.
type Storage interface {
	Connect(*config.StorageConnection) error
	Close() error
	Create(storage.Event) error
	Update(int, storage.Event) error
	Delete(int) error
	ListDay(time.Time, int) ([]storage.Event, error)
	ListWeek(time.Time, int) ([]storage.Event, error)
	ListMonth(time.Time, int) ([]storage.Event, error)
	Clean() error
	IsBusy(time.Time, time.Time, int) (bool, error)
}

// New app.
func New(logger Logger, storage Storage) *App {
	app := App{
		logger,
		storage,
	}
	return &app
}

func (a *App) StorageConnect(storageConfig *config.StorageConnection) error {
	err := a.Storage.Connect(storageConfig)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) StorageClose() error {
	err := a.Storage.Close()
	if err != nil {
		return err
	}

	return nil
}

// Create event.
func (a *App) CreateEvent(userID int, title, description string, startDate, endDate time.Time) error {
	// Check parameters.
	if userID == 0 {
		return ErrNoUserID
	}

	if title == "" {
		return ErrEmptyTitle
	}

	if startDate.IsZero() {
		return ErrNoStartDate
	}

	if endDate.IsZero() {
		return ErrNoEndDate
	}

	if endDate.Before(startDate) {
		return ErrEndBeforeStart
	}

	// Check if busy.
	busy, err := a.Storage.IsBusy(startDate, endDate, userID)
	if err != nil {
		return err
	}

	if busy {
		return ErrTimeBusy
	}

	// Create event
	err = a.Storage.Create(storage.Event{
		Title:       title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		UserID:      userID,
	})
	if err != nil {
		return err
	}

	return nil
}

// Update event.
func (a *App) UpdateEvent(id int, title, description string, startDate, endDate time.Time, userID int) error {
	// Check parameters.
	if id == 0 {
		return ErrNoEventID
	}

	if userID == 0 {
		return ErrNoUserID
	}

	if title == "" {
		return ErrEmptyTitle
	}

	if startDate.IsZero() {
		return ErrNoStartDate
	}

	if endDate.IsZero() {
		return ErrNoEndDate
	}

	if endDate.Before(startDate) {
		return ErrEndBeforeStart
	}

	// Check if busy.
	busy, err := a.Storage.IsBusy(startDate, endDate, userID)
	if err != nil {
		return err
	}

	if busy {
		return ErrTimeBusy
	}

	err = a.Storage.Update(id, storage.Event{
		ID:          id,
		Title:       title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		UserID:      userID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteEvent(id int) error {
	// Check parameters.
	if id == 0 {
		return ErrNoEventID
	}

	err := a.Storage.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ListEventsForDay(startDate time.Time, userID int) ([]storage.Event, error) {
	// Check parameters.
	if userID == 0 {
		return nil, ErrNoUserID
	}

	if startDate.IsZero() {
		return nil, ErrNoStartDate
	}

	// Get events for day.
	events, err := a.Storage.ListDay(startDate, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (a *App) ListEventsForWeek(startDate time.Time, userID int) ([]storage.Event, error) {
	// Check parameters.
	if userID == 0 {
		return nil, ErrNoUserID
	}

	if startDate.IsZero() {
		return nil, ErrNoStartDate
	}

	// Get events for week.
	events, err := a.Storage.ListWeek(startDate, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (a *App) ListEventsForMonth(startDate time.Time, userID int) ([]storage.Event, error) {
	// Check parameters.
	if userID == 0 {
		return nil, ErrNoUserID
	}

	if startDate.IsZero() {
		return nil, ErrNoStartDate
	}

	// Get events for month.
	events, err := a.Storage.ListMonth(startDate, userID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (a *App) Clean() error {
	err := a.Storage.Clean()
	if err != nil {
		return err
	}

	return nil
}
