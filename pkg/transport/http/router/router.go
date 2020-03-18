package router

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type (
	Module struct {
		Router func(e *echo.Echo)
	}
)

// Params are the parameters required to build the router
type Params struct {
	fx.In

	Modules     []Module              `group:"http"`
	Middlewares []echo.MiddlewareFunc `group:"middleware"`
}

// New creates a new instance of a *mux.Router with all of the modules added
func NewRouter(params Params) *echo.Echo {
	e := echo.New()
	for _, module := range params.Modules {
		module.Router(e)
	}

	e.Use(params.Middlewares...)

	return e
}
