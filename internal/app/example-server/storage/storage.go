package storage

import (
	"example-server/internal/app/example-server/storage/cache"
	"example-server/internal/app/example-server/storage/db"
	"github.com/hongjundu/go-color-logger"
)

func Init() error {
	logger.Debug("[storage] Init")

	if err := db.Init(); err != nil {
		return err
	}

	if err := cache.Init(); err != nil {
		return err
	}

	return nil
}
