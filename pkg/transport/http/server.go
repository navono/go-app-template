package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-app-template/internal/dependency"
	"go-app-template/pkg/transport/http/router"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"
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
		router.NewRouter,
	),
	Constructor: func(e *echo.Echo, logger *zap.Logger, getter dependency.ConfigGetter) WebServer {
		//e.Debug = true
		e.Use(echoMw.Recover())

		e.Server.Addr = fmt.Sprintf("%s:%d", getter.GetString("http-host"), getter.GetInt("http-port"))

		return &server{e, logger}
	},
	InvokeFunc: Invoke,
}

type (
	WebServer interface {
		ListenAndServe() error
		Shutdown(ctx context.Context) error
	}

	server struct {
		echo   *echo.Echo
		logger *zap.Logger
	}
)

func (ws *server) ListenAndServe() error {
	ws.logger.Info("Starting HTTP Server", zap.String("Addr", ws.echo.Server.Addr))
	if err := graceful.ListenAndServe(ws.echo.Server, 5*time.Second); err != nil {
		ws.logger.Error("Could not start Server", zap.Error(err))
		return err
	}
	return nil
}

func (ws *server) Shutdown(ctx context.Context) error {
	ws.logger.Info("Stopping HTTP Server")
	return ws.echo.Shutdown(ctx)
}

// Params are the dependencies required to start the http
type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Server    WebServer
	Logger    *zap.Logger
}

// Invoke is the function that is called to start the http
func Invoke(params Params) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: StartServer(params.Server, params.Logger),
		OnStop:  StopServer(params.Server, params.Logger),
	})
}

// StartServer creates a closure that will start the http when called
func StartServer(s WebServer, logger *zap.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		go func() {
			if err := s.ListenAndServe(); err != nil {
				logger.Error("Could not start Server", zap.Error(err))
			}
		}()
		return nil
	}
}

func StopServer(s WebServer, logger *zap.Logger) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if err := s.Shutdown(context.Background()); err != nil {
			logger.Error("Error when shutting down Server", zap.Error(err))
			return fmt.Errorf("error shutting down Server (%w)", err)
		}
		return nil
	}
}
