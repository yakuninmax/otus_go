package app

import (
	"context"
	"time"

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

// Storage interface.
type Storage interface {
	Connect(context.Context, string) error
	Close(context.Context)
	Create(context.Context, storage.Event) (int, error)
	Update(context.Context, int, storage.Event) error
	Delete(context.Context, int) error
	ListDay(context.Context, time.Time) ([]storage.Event, error)
	ListWeek(context.Context, time.Time) ([]storage.Event, error)
	ListMonth(context.Context, time.Time) ([]storage.Event, error)
	Clean(context.Context) error
}

// New app.
func New(logger Logger, storage Storage) *App {
	app := App{
		logger,
		storage,
	}
	return &app
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
