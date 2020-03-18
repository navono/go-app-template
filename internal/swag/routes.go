package swag

import (
	"net/http"

	"go-app-template/pkg/transport/http/router"

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
		Method:      http.MethodGet,
		Path:        "/swagger/*",
		HandlerFunc: echoSwagger.WrapHandler,
	}
}
