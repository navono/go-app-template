package serve

import (
	"go-app-template/internal/config"
	"go-app-template/internal/dependency"
	"go-app-template/internal/logging"
	"go-app-template/internal/middleware"
	"go-app-template/internal/response"
	"go-app-template/pkg/service/hello"
	"go-app-template/pkg/transport/http"

	"github.com/spf13/cobra"
)

// NewCommand creates an instance of the Serve command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web-http",
		Long:  "Start the example web-sever",
	}
	cmd.Run = Serve(newWebApplicationBuilder(cmd))
	return cmd
}

// Serve produces the function that is called when the
// command is called
func Serve(builder dependency.Builder) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		builder.
			WithModule(middleware.Module).
			WithModule(hello.Module).
			Build().
			Run()
	}
}

func newWebApplicationBuilder(cmd *cobra.Command) dependency.Builder {
	return dependency.NewBuilder(cmd).
		WithService(config.Service).
		WithService(logging.Service).
		WithService(response.Service).
		WithService(http.Service)
}
