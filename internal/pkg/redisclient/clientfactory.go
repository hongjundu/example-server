package redisclient

import (
	"example-server/internal/pkg/env"
	"github.com/go-redis/redis"
	"github.com/hongjundu/go-color-logger"
	"strings"
)

const (
	redisModeStandalone = "standalone"
	redisModeCluster    = "cluster"
	redisModeSentinel   = "sentinel"
)

type RedisClient interface {
	Ping() *redis.StatusCmd
}

func NewRedisClient() (client RedisClient, err error) {
	logger.Debug("[redisclient] NewRedisClient")

	if len(env.Env.RedisHost) == 0 {
		logger.Fatal("[redisclient] No redis address was configed")
	}

	if strings.Compare(env.Env.RedisMode, redisModeStandalone) == 0 {
		client = redis.NewClient(&redis.Options{
			Addr:     env.Env.RedisHost,
			Password: env.Env.RedisPassword, // redis password
			DB:       env.Env.RedisDb,       // use default DB
		})

	} else if strings.Compare(env.Env.RedisMode, redisModeSentinel) == 0 {
		if len(env.Env.RedisMasterName) == 0 {
			logger.Fatal("[redisclient] No master name was configed")
		}

		redisAddrs := strings.Split(env.Env.RedisHost, ",")

		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    env.Env.RedisMasterName,
			SentinelAddrs: redisAddrs,
			Password:      env.Env.RedisPassword,
			DB:            env.Env.RedisDb,
		})

	} else if strings.Compare(env.Env.RedisMode, redisModeCluster) == 0 {
		redisAddrs := strings.Split(env.Env.RedisHost, ",")

		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisAddrs,
			Password: env.Env.RedisPassword,
		})

	} else {
		logger.Fatal("[redisclient] unsupported redis mode: %s", env.Env.RedisMode)
	}

	if pong, e := client.Ping().Result(); e == nil {
		logger.Info("[redisclient]", "pong", pong)
	} else {
		err = e
		logger.Error("[redisclient] NewRedisClient", "error", err)
		return
	}

	return
}
