package router_test

import (
	"context"
	"net/http"
	"reflect"
	"sort"
	"testing"

	"go-app-template/internal/config"
	"go-app-template/internal/dependency"
	"go-app-template/internal/logging"
	"go-app-template/internal/response"
	"go-app-template/pkg/transport/http/router"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

type PrinterFunc func(string, ...interface{})

func (p PrinterFunc) Printf(string, ...interface{}) {}

func TestNew(t *testing.T) {
	module := router.Module{
		//Path:   "test",
		//Router: func(router *mux.Router) {},
		Router: func(e *echo.Echo) {},
	}
	cmd := &cobra.Command{}
	annotated := fx.Annotated{
		Group: "http",
		Target: func() router.Module {
			return module
		},
	}

	module2 := router.Module{
		//Path:   "/test2",
		//Router: func(router *mux.Router) {},
		Router: func(e *echo.Echo) {},
	}

	annotated2 := fx.Annotated{
		Group: "http",
		Target: func() router.Module {
			return module2
		},
	}

	expectedRoutes := []string{"/test", "/test2"}
	gotRoutes := make([]string, 0, len(expectedRoutes))
	builder := dependency.NewBuilder(cmd).
		WithService(config.Service).
		WithService(logging.Service).
		WithService(response.Service).
		//WithService(router.Service).
		WithModule(fx.Provide(annotated, annotated2)).
		WithModule(fx.Logger(PrinterFunc(func(format string, args ...interface{}) {
			t.Logf(format, args...)
		}))).
		WithInvoke(func(handler http.Handler) {
			muxRouter, ok := handler.(*mux.Router)
			if !ok {
				t.Fatal("cannot convert handler to *mux.Router")
			}
			err := muxRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
				routeString, err := route.GetPathTemplate()
				if err != nil {
					t.Errorf("got error (%s), when getting route", err)
				}
				gotRoutes = append(gotRoutes, routeString)
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	app := builder.BuildTest(t)
	go app.Run()
	if err := app.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}

	sort.Strings(expectedRoutes)
	sort.Strings(gotRoutes)
	if !reflect.DeepEqual(expectedRoutes, gotRoutes) {
		t.Fatalf("expected routes to be (%+v), got (%+v)", expectedRoutes, gotRoutes)
	}
}
