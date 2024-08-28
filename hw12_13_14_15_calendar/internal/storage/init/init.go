package initstorage

import (
	"fmt"

	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/yakuninmax/otus_go/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(config *config.StorageConfig) (app.Storage, error) {
	var storage app.Storage

	switch config.Type {
	case "in-memory":
		storage = memorystorage.New()
	case "sql":
		storage = sqlstorage.New()
	default:
		return nil, fmt.Errorf("unsupported storage type: %s" + config.Type)
	}

	return storage, nil
}
