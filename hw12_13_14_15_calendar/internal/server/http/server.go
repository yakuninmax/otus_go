package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	config *config.HTTPConfig
	logger Logger
	app    Application
	server *http.Server
}

type Logger interface {
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debug(...interface{})
}

type Application interface {
	StorageConnect(*config.StorageConnection) error
	StorageClose() error
	CreateEvent(int, string, string, time.Time, time.Time) error
	UpdateEvent(int, string, string, time.Time, time.Time, int) error
	DeleteEvent(int) error
	ListEventsForDay(time.Time, int) ([]storage.Event, error)
	ListEventsForWeek(time.Time, int) ([]storage.Event, error)
	ListEventsForMonth(time.Time, int) ([]storage.Event, error)
	Clean() error
}

func New(config *config.HTTPConfig, logger Logger, app Application) *Server {
	return &Server{
		config: config,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context, logger Logger) error {
	logFile, err := os.OpenFile(s.config.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o664)
	if err != nil {
		log.Fatal(err)
	}

	// Create new router.
	mux := http.NewServeMux()

	// Configure router.
	hello := http.HandlerFunc(HelloHandler)
	mux.Handle("/", handlers.LoggingHandler(logFile, hello))

	// Configure server.
	s.server = &http.Server{
		Addr:         net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port)),
		Handler:      mux,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Run server.
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	return nil
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
