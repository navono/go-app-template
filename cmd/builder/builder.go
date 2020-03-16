package builder

import (
	"go-app-template/internal/config"
	"go-app-template/internal/dependency"
	"go-app-template/internal/logging"
	"go-app-template/internal/response"
	"go-app-template/pkg/transport/http"

	"github.com/spf13/cobra"
)

func NewApplicationBuilder(cmd *cobra.Command) dependency.Builder {
	return dependency.NewBuilder(cmd).
		WithService(config.Service).
		WithService(logging.Service).
		WithService(response.Service).
		WithService(http.Service)
}
