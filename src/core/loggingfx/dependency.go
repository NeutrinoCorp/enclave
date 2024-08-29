package loggingfx

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/application"
	"github.com/neutrinocorp/geck/logging"
)

var StdLoggerModule = fx.Module("logger_std",
	fx.Provide(
		func() *log.Logger {
			return log.New(os.Stdout, "", 0)
		},
		fx.Annotate(
			logging.NewStdLogger,
			fx.As(new(logging.Logger)),
		),
	),
)

var ZerologLoggerModule = fx.Module("logger_zerolog",
	fx.Provide(
		func() zerolog.Logger {
			return zerolog.New(os.Stdout).With().Timestamp().Logger()
		},
		fx.Annotate(
			logging.NewZerologLogger,
			fx.As(new(logging.Logger)),
		),
	),
)

var ZerologAppLoggerModule = fx.Module("logger_zerolog_app",
	fx.Provide(
		func(cfg application.Config) zerolog.Logger {
			return logging.NewApplicationLogger(cfg, os.Stdout)
		},
		fx.Annotate(
			logging.NewZerologLogger,
			fx.As(new(logging.Logger)),
		),
	),
)

func DecorateLoggerWithModule(moduleName string) any {
	return func(logger logging.Logger) logging.Logger {
		return logger.Module(moduleName)
	}
}
