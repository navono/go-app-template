package swagger

import (
	"go-app-template/pkg/transport/http/router"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotated{
		Group:  "http",
		Target: RegisterSwagHandler,
	},
)

func RegisterSwagHandler() router.Module {
	return router.Module{
		Router: func(e *echo.Echo) {
			e.GET("/swagger/*", echoSwagger.WrapHandler)
		},
	}
}
