package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/actuatorfx"
	"github.com/neutrinocorp/geck/applicationfx"
	"github.com/neutrinocorp/geck/observability/loggingfx"
	"github.com/neutrinocorp/geck/securityfx"
	"github.com/neutrinocorp/geck/transportfx"
	"github.com/neutrinocorp/geck/validationfx"

	"http-server/pingfx"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("error loading .env file, using system env variables")
	}

	app := fx.New(
		applicationfx.ApplicationModule,
		loggingfx.ZerologAppLoggerModule,
		actuatorfx.ActuatorModule,
		validationfx.GoPlaygroundValidationModule,
		securityfx.CognitoModule,
		transportfx.TransportModuleHTTP,
		transportfx.TransportJWTModuleHTTP,
		pingfx.PingModule,
	)
	app.Run()
}
