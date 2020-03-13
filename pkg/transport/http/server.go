package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-app-template/internal/dependency"
	"go-app-template/pkg/transport/http/router"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Service allows the service to be used in the dependency builder
var Service = dependency.Service{
	ConfigFunc: func(flags dependency.FlagSet) {
		flags.String("http-host", "127.0.0.1", "The IP to start on")
		flags.Int("http-port", 8080, "The port to start the web Server on")
		flags.Duration("read-timeout", 10*time.Second, "The read timeout for the HTTP Server")
		flags.Duration("read-header-timeout", 20*time.Second, "The read header timeout for the HTTP Server")
		flags.Duration("write-timeout", 20*time.Second, "The write timeout for the HTTP Server")
		flags.Duration("idle-timeout", 10*time.Second, "The idle timeout for the HTTP Server")
		flags.Int("max-header-bytes", http.DefaultMaxHeaderBytes, "The maximum size that the HTTP header can be in bytes")
	},
	Dependencies: fx.Provide(
		router.NewRouter, NewServer,
	),
	InvokeFunc: Invoke,
	Constructor: func(server *http.Server) Server {
		return server
	},
}

// New creates a new instance of the *http.Server configured by the config
// you decided
func NewServer(router http.Handler, getter dependency.ConfigGetter) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", getter.GetString("http-host"), getter.GetInt("http-port")),
		Handler:           router,
		ReadTimeout:       getter.GetDuration("read-timeout"),
		ReadHeaderTimeout: getter.GetDuration("read-header-timeout"),
		WriteTimeout:      getter.GetDuration("write-timeout"),
		IdleTimeout:       getter.GetDuration("idle-timeout"),
		MaxHeaderBytes:    getter.GetInt("max-header-bytes"),
	}
}

// Params are the dependencies required to start the http
type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Server    Server
	Logger    *zap.Logger
}

// Invoke is the function that is called to start the http
func Invoke(params Params) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: StartServer(params.Server, params.Logger),
		OnStop:  StopServer(params.Server, params.Logger),
	})
}

// Server is an interface that abstracts the *http.Server
type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// StartServer creates a closure that will start the http when called
func StartServer(server Server, logger *zap.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		logger.Info("Starting HTTP Server")
		go func() {
			if err := server.ListenAndServe(); err != nil {
				logger.Error("Could not start Server", zap.Error(err))
			}
		}()
		return nil
	}
}

// StopServer creates a closure that will stop the http
func StopServer(server Server, logger *zap.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		logger.Info("Stopping HTTP Server")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Error when shutting down Server")
			return fmt.Errorf("error shutting down Server (%w)", err)
		}
		return nil
	}
}
