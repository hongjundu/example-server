package server

import (
	"context"
	"example-server/internal/app/example-server/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-level-logger"
	"github.com/rs/cors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	return &Server{
		router: gin.New(),
	}
}

func myLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		logger.Debugf("Method: %s URL: %s", c.Request.Method, c.Request.URL)

		c.Next()

		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()

		logger.Debugf("Latency: %v Status: %v", latency, status)
	}
}

func (server *Server) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO: verify token

		//c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Unauthorized")))
		//c.Abort()

		c.Next()
	}
}

func (server *Server) configRouter() {
	logger.Debugf("[Server] configRouter")

	server.router.Use(gin.Logger())
	//server.router.Use(myLogger())
	server.router.Use(gin.Recovery())

	server.router.NoRoute(ginHandlerFunc(server.notFoundHandler))

	v1 := server.router.Group("/v1")
	v1.POST("/login", ginHandlerFunc(server.loginHandler))
	v1.GET("/version", ginHandlerFunc(server.versionHandler))
	v1.Use(server.AuthRequired())
	{
		v1.POST("/tasks", ginHandlerFunc(server.createTaskHandler))
	}

}

func (server *Server) Run(port int) error {
	logger.Debugf("[Server] Run")

	server.configRouter()

	// CORS
	handler := cors.Default().Handler(server.router)
	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"authorization", "Content-Type"},
		OptionsPassthrough: true,
		AllowCredentials:   true,
	})
	handler = c.Handler(handler)

	if port <= 0 {
		port = 8000
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	if err := storage.Init(); err != nil {
		logger.Fatalf("[Server] storage.Init(): %+v", err)
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("[Server] listen error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("[Server] shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("[Server] shutdown error: %v", err)
	}

	logger.Infof("[Server] server exit.")

	return nil
}
