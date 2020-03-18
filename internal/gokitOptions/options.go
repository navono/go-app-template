package gokitOptions

import (
	kitZap "github.com/go-kit/kit/log/zap"
	kitTransport "github.com/go-kit/kit/transport"
	kitHttp "github.com/go-kit/kit/transport/http"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(
	NewKitOption,
)

func NewKitOption(logger *zap.Logger) []kitHttp.ServerOption {
	kitLogger := kitZap.NewZapSugarLogger(logger, zap.ErrorLevel)
	return []kitHttp.ServerOption{
		kitHttp.ServerErrorHandler(kitTransport.NewLogErrorHandler(kitLogger)),
		kitHttp.ServerErrorEncoder(kitHttp.DefaultErrorEncoder),
	}
}
