package hello

import (
	"net/http"

	"go-app-template/pkg/transport/http/router"

	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHelloService,
	fx.Annotated{
		Group:  "http",
		Target: RegisterHandler,
	},
)

type HandlerParams struct {
	fx.In

	Options      []kitHttp.ServerOption
	HelloService *HelloService
}

func RegisterHandler(params HandlerParams) router.Module {
	return router.Module{
		Method: http.MethodGet,
		Path:   "/hello",
		// Also, we can use normal request handler
		HandlerFunc: echo.WrapHandler(
			kitHttp.NewServer(
				makeHelloEndpoint(params.HelloService),
				decodeHelloRequest,
				encodeResponse,
				params.Options...)),
	}
}
