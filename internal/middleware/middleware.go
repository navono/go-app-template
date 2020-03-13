package middleware

import (
	"go-app-template/internal/logging"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

// Module allows the default middlewares to be registered to an app
var Module = fx.Provide(
	fx.Annotated{
		Group: "middleware",
		Target: func(logger logging.PrintLogger) mux.MiddlewareFunc {
			return handlers.RecoveryHandler(
				handlers.RecoveryLogger(logger),
			)
		},
	},
	fx.Annotated{
		Group: "middleware",
		Target: func() mux.MiddlewareFunc {
			return gziphandler.GzipHandler
		},
	},
)
