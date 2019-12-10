package main

import (
	"example-server/internal/app/example-server/server"
	"example-server/internal/pkg/env"
	"github.com/hongjundu/go-level-logger"
	"github.com/kelseyhightower/envconfig"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /v1

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
