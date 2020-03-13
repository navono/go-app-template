package hello

import (
	"go-app-template/pkg/transport/http/router"

	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewHelloHandler,
	fx.Annotated{
		Group:  "http",
		Target: RegisterHandler,
	},
)

type HandlerParams struct {
	fx.In

	HelloHandler *HelloHandler
}

func RegisterHandler(params HandlerParams) router.Module {
	return router.Module{
		Path: "hello",
		Router: func(router *mux.Router) {
			//router.Handle("/", params.HelloHandler.Get)
			//router.HandleFunc("/", params.HelloHandler.Get2)

			router.Handle("/", kitHttp.NewServer(
				makeHelloEndpoint(params.HelloHandler),
				decodeHelloRequest,
				encodeResponse))
		},
	}
}
