package serve

import (
	"go-app-template/cmd/builder"
	"go-app-template/internal/dependency"
	"go-app-template/internal/middleware"
	"go-app-template/pkg/service/hello"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

// NewCommand creates an instance of the Serve command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the web-http",
		Long:  "Start the example web-sever",
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")

	cmd.Run = Serve(builder.NewApplicationBuilder(cmd))
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
