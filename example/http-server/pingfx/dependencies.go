package pingfx

import (
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/loggingfx"
	"github.com/neutrinocorp/geck/transportfx"
	"http-server/ping"
)

var PingModule = fx.Module("ping",
	fx.Decorate(
		loggingfx.DecorateLoggerWithModule("ping"),
	),
	fx.Provide(
		transportfx.AsVersionedControllerHTTP(ping.NewControllerHTTP),
	),
)
