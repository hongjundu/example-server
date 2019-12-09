package main

import (
	"example-server/internal/app/example-server/server"
	"example-server/internal/pkg/env"
	"github.com/hongjundu/go-level-logger"
	"github.com/kelseyhightower/envconfig"
)

func main() {

	if err := envconfig.Process("example_server", &env.Env); err != nil {
		logger.Fatal(err)
	}

	logger.Init(env.Env.LogLevel, "example-server", env.Env.LogPath, 100, 3, 30)
	logger.Debugf("[main] env: %+v", env.Env)

	s := server.NewServer()
	if e := s.Run(env.Env.Port); e != nil {
		logger.Fatalf("%v", e)
	}
}
