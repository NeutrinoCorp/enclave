package loggingfx

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/application"
	logging2 "github.com/neutrinocorp/geck/observability/logging"
)

var StdLoggerModule = fx.Module("logger_std",
	fx.Provide(
		func() *log.Logger {
			return log.New(os.Stdout, "", 0)
		},
		fx.Annotate(
			logging2.NewStdLogger,
			fx.As(new(logging2.Logger)),
		),
	),
)

var ZerologLoggerModule = fx.Module("logger_zerolog",
	fx.Provide(
		func() zerolog.Logger {
			return zerolog.New(os.Stdout).With().Timestamp().Logger()
		},
		fx.Annotate(
			logging2.NewZerologLogger,
			fx.As(new(logging2.Logger)),
		),
	),
)

var ZerologAppLoggerModule = fx.Module("logger_zerolog_app",
	fx.Provide(
		func(cfg application.Config) zerolog.Logger {
			return logging2.NewApplicationLogger(cfg, os.Stdout)
		},
		fx.Annotate(
			logging2.NewZerologLogger,
			fx.As(new(logging2.Logger)),
		),
	),
)

func DecorateLoggerWithModule(moduleName string) any {
	return func(logger logging2.Logger) logging2.Logger {
		return logger.Module(moduleName)
	}
}
