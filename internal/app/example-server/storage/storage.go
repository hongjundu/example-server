package storage

import (
	"example-server/internal/app/example-server/storage/cache"
	"example-server/internal/app/example-server/storage/db"
	"github.com/hongjundu/go-level-logger"
)

func Init() error {
	logger.Debugf("[storage] Init")

	if err := db.Init(); err != nil {
		return err
	}

	if err := cache.Init(); err != nil {
		return err
	}

	return nil
}
