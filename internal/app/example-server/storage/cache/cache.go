package cache

import (
	"example-server/internal/pkg/redisclient"
	"github.com/hongjundu/go-color-logger"
)

var redisClient redisclient.RedisClient

func Init() error {
	var err error
	if redisClient, err = redisclient.NewRedisClient(); err != nil {
		logger.Error("[cache] create redis client", "error", err)
	}
	return nil
}
