package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go-app-template/internal/config"
	"go-app-template/internal/dependency"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Service is the exported variable that can be used by the framework package
var Service = dependency.Service{
	Dependencies: fx.Provide(
		NewPrintLogger,
		fx.Annotated{
			Group:  "middleware",
			Target: NewMidlleware,
		},
		NewLoggerFactory,
	),
	ConfigFunc: func(set dependency.FlagSet) {
		set.String("app-name", filepath.Base(os.Args[0]), "The name of the application being configured")
		set.String("app-version", "dev", "The version of the application being configured")
		set.String("environment", "test", "The environment that the application is deployed in")
		set.String("logger", "development", "Whether to log in development mode.")
		set.StringSlice("excluded-headers", []string{"Authorization"}, "Which headers to hide from the request log")
	},
	Constructor: func(lf LoggerFactory, settings dependency.ConfigGetter) (*zap.Logger, error) {
		logger, err := lf.Logger(settings)
		if err != nil {
			return nil, err
		}

		logger.Info("Process environment", zap.String("app-version", settings.GetString("app-version")),
			zap.Int("pid", os.Getpid()))

		return logger, nil
	},
}

// LoggerConstructor is a type that can give you an instance of a logger
type LoggerConstructor func(options ...zap.Option) (*zap.Logger, error)

// NewLoggerFactory will create a new instance of a logger factory
func NewLoggerFactory(fileConfig *config.FileConfig) LoggerFactory {
	var production func(options ...zap.Option) (*zap.Logger, error)
	if fileConfig == nil {
		production = zap.NewProduction
	} else {
		logConfig := fileConfig.Log
		hook := &lumberjack.Logger{
			Filename:   logConfig.Filename,    //filePath
			MaxSize:    logConfig.FileMaxSize, // megabytes
			MaxBackups: 10000,
			MaxAge:     logConfig.FileMaxAge, //days
			Compress:   logConfig.Compress,   // disabled by default
		}
		enConfig := zap.NewProductionEncoderConfig() //生成配置
		enConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		w := zapcore.AddSync(hook)
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(enConfig), //编码器配置
			w,                                   //打印到控制台和文件
			getZapLevel(logConfig.FileLevel),    //日志等级
		)
		zap.NewProductionConfig()
		production = func(options ...zap.Option) (logger *zap.Logger, err error) {
			return zap.New(core, zap.AddCallerSkip(logConfig.Skip), zap.AddCaller()).WithOptions(options...), nil
		}
	}

	return LoggerFactory{
		LoggerConstructors: map[string]LoggerConstructor{
			"production":  production,
			"development": zap.NewDevelopment,
			"nop": func(options ...zap.Option) (logger *zap.Logger, err error) {
				return zap.NewNop(), nil
			},
		},
	}
}

// LoggerFactory is a type that can create instances of loggers
type LoggerFactory struct {
	LoggerConstructors map[string]LoggerConstructor
}

// Logger creates a new instance of a *zap.Logger
func (f LoggerFactory) Logger(settings dependency.ConfigGetter) (*zap.Logger, error) {
	options := []zap.Option{
		zap.Fields(
			zap.String("app-name", settings.GetString("app-name")),
			zap.String("app-version", settings.GetString("app-version")),
			zap.String("environment", settings.GetString("environment")),
		),
	}

	loggerType := settings.GetString("logger")
	loggerConstructor, ok := f.LoggerConstructors[loggerType]
	if !ok {
		return nil, fmt.Errorf("the logger type (%s), is not a valid logger", loggerType)
	}

	logger, err := loggerConstructor(options...)
	if err != nil {
		return nil, fmt.Errorf("could not create instance of logger, got error (%w)", err)
	}

	return logger, nil
}

// NewPrintLogger creates a new instance of the PrintLogger
func NewPrintLogger(logger *zap.Logger) PrintLogger {
	return PrintLogger{
		Logger: logger,
	}
}

// PrintLogger implements a generalised logging interface
type PrintLogger struct {
	Logger *zap.Logger
}

// Println prints the arguments to the zap logger
func (p PrintLogger) Println(arguments ...interface{}) {
	p.Logger.Info(fmt.Sprintln(arguments...))
}

func getZapLevel(s string) zapcore.Level {
	switch strings.ToLower(s) {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
