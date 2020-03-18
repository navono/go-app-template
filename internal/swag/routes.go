package swag

import (
	"go-app-template/internal/dependency"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"go-app-template/pkg/transport/http/router"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotated{
		Group:  "http",
		Target: RegisterSwagHandler,
	},
)

func RegisterSwagHandler(getter dependency.ConfigGetter) router.Module {
	//host := getter.GetString("http-host")
	//port := getter.GetInt("http-port")
	//url := fmt.Sprintf("%s:%d", host, port)
	//swagPath := fmt.Sprintf("%s/api/swagger/doc.json", url)

	return router.Module{
		Path: "swagger",
		Router: func(router *mux.Router) {
			router.Handle("/", handlers.MethodHandler{
				http.MethodGet: httpSwagger.WrapHandler,
			})
		},
	}
}
