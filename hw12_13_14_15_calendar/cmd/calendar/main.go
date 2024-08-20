package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage/memory"
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
		os.Stdout.WriteString("invalid config file path: " + err.Error())
		os.Exit(1) //nolint:gocritic
	}

	// Read config file.
	config, err := config.NewConfig(configFile)
	if err != nil {
		os.Stdout.WriteString("failed to read config file: " + err.Error())
		os.Exit(1) //nolint:gocritic
	}

	// Create logger.
	logg := logger.New(config.Logger.Level)

	// Create storage.
	storage := memorystorage.New()

	// Create calendar app.
	calendar := app.New(logg, storage)

	// Create HTTP server.
	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
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
