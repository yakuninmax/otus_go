package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/server/http"
	storage "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage/init"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

func main() {
	// Parse input flags.
	flag.Parse()

	// Print version info.
	if flag.Arg(0) == "version" {
		printVersion()
		os.Exit(0)
	}

	// Read config file.
	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Create logger.
	logger, err := logger.New(&config.LoggerConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create storage.
	storage, err := storage.New(&config.StorageConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Connect db.
	err = storage.Connect(&config.StorageConnection)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer storage.Close()

	// Create calendar app.
	calendar := app.New(logger, storage)

	// Create HTTP server.
	server := internalhttp.New(&config.HTTPConfig, logger, calendar)

	// Create context.
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("calendar is running...")

	if err := server.Start(ctx, logger); err != nil {
		log.Fatalln(err.Error()) //nolint:gocritic
	}
}
