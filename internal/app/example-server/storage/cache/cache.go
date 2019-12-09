package cache

import (
	"example-server/internal/pkg/redisclient"
	"github.com/hongjundu/go-level-logger"
)

var redisClient redisclient.RedisClient

func Init() error {
	var err error
	if redisClient, err = redisclient.NewRedisClient(); err != nil {
		logger.Errorln("[cache] create redis client failed")
	}
	return nil
}
