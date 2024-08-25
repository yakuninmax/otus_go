package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	// Parse input flags.
	flag.Parse()

	// Print version info.
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Validate config file path.
	if err := validateConfigPath(configFile); err != nil {
		log.Fatalf("invalid config file path: %s", err.Error())
	}

	// Read config file.
	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	// Create logger.
	logger, err := logger.New(config.LoggerConfig.Level, config.LoggerConfig.Colors, config.LoggerConfig.FullTimestamp)
	if err != nil {
		log.Fatalf("failed to configure logger: %s", err.Error())
	}

	// Create context.
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Create storage.
	var storage app.Storage

	switch config.Type {
	case "in-memory":
		storage = memorystorage.New()
	case "sql":
		storage = sqlstorage.New()
		// Get connection string.
		connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.StorageConnection.Host,
			config.StorageConnection.Port,
			config.StorageConnection.User,
			config.StorageConnection.Password,
			config.StorageConnection.Database)

		err := storage.Connect(ctx, connectionString)
		if err != nil {
			logger.Error("database connection failed: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	default:
		logger.Error("unsupported storage type: " + config.Type)
		cancel()
		os.Exit(1)
	}

	// Create calendar app.
	calendar := app.New(logger, storage)

	// Create HTTP server.
	server := internalhttp.NewServer(logger, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}

func validateConfigPath(configFilePath string) error {
	s, err := os.Stat(configFilePath)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory", configFilePath)
	}
	return nil
}
