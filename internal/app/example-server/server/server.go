package server

import (
	"context"
	"crypto/rsa"
	"example-server/internal/app/example-server/storage"
	_ "example-server/swagger"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-color-logger"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	srv           *http.Server
	router        *gin.Engine
	jwtPublicKey  *rsa.PublicKey
	jwtPrivateKey *rsa.PrivateKey
	enforcer      *casbin.Enforcer
	bigCache      *bigcache.BigCache
	exit          chan struct{}
}

func NewServer() *Server {
	return &Server{
		router: gin.New(),
		exit:   make(chan struct{}),
	}
}

func myLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		logger.Debug("Request", "Method", c.Request.Method, "URL", c.Request.URL)

		c.Next()

		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()

		logger.Debug("Request", "latency", latency, "status", status)
	}
}

func (server *Server) configRouter() {
	logger.Debug("[Server] configRouter")

	server.router.Use(gin.Logger())
	//server.router.Use(myLogger())
	server.router.Use(gin.Recovery())

	server.router.NoRoute(server.ginHandlerFunc(server.notFoundHandler))

	v1 := server.router.Group("/v1")
	v1.POST("/login", server.ginHandlerFunc(server.loginHandler))
	v1.GET("/version", server.ginHandlerFunc(server.versionHandler))
	v1.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "EXAMPLE_SERVER_DISABLE_SWAGGER"))
	v1.Use(server.jwtTokenRequired())
	{
		v1.POST("/logout", server.ginHandlerFunc(server.logoutHandler))
		v1.GET("/hello", server.acl("data", "read"), server.ginHandlerFunc(server.readHandler))
		v1.POST("/hello", server.acl("data", "write"), server.ginHandlerFunc(server.writeHandler))
	}
}

func (server *Server) Run(port int) error {
	logger.Debug("[Server] Run")

	if err := server.loadJwtKeys(); err != nil {
		logger.Fatal("[Server] loadJwtKeys", "error", err)
	}

	if err := server.createEnforcer(); err != nil {
		logger.Fatal("[Server] createEnforcer", "error", err)
	}

	if bigCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(1 * time.Minute)); err != nil {
		logger.Fatal("[Server] bigcache.NewBigCache", "error", err)
	} else {
		server.bigCache = bigCache
	}

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

	server.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	if err := storage.Init(); err != nil {
		logger.Fatal("[Server] storage.Init()", "error", err)
	}

	go func() {
		// service connections
		if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("[Server] listen", "error", err)
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

	server.Close()

	return nil
}

func (server *Server) IsExiting() bool {
	select {
	case <-server.exit:
		return true
	default:
		return false
	}
}

func (server *Server) Close() {
	logger.Debug("[Server] Close")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()

	// Indicate that we do not accept request anymore
	close(server.exit)

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown", "error", err)
	}
	// Optionally, you could run svr.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("[HttpServer] shutdown gracefully")
}
