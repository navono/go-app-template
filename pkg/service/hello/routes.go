package hello

import (
	"go-app-template/pkg/transport/http/router"

	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
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

	HelloService *HelloService
}

func RegisterHandler(params HandlerParams) router.Module {
	return router.Module{
		Path: "hello",
		Router: func(router *mux.Router) {
			router.Handle("/", kitHttp.NewServer(
				makeHelloEndpoint(params.HelloService),
				decodeHelloRequest,
				encodeResponse))
		},
	}
}
