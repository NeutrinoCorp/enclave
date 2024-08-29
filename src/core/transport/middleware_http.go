package transport

import (
	"errors"
	"strconv"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/logging"
	"github.com/neutrinocorp/geck/security"
)

// Error handler

var _ echo.HTTPErrorHandler = HandleEchoError

func HandleEchoError(srcErr error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	errs := convertContainerErrorsEcho(srcErr)
	_ = c.JSON(errs.Code, errs)
}

// General-purposed middlewares

type DefaultEchoMiddlewareParams struct {
	fx.In

	Config ConfigHTTP
	Logger logging.Logger
}

// NewDefaultEchoMiddlewareGroup allocates an Echo middleware group. An array is returned to guarantee
// ordering within this group. Any other middlewares outside this group will be injected into application in
// a non-deterministic way (regarding ordering).
func NewDefaultEchoMiddlewareGroup(params DefaultEchoMiddlewareParams) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.RequestID(),
		LogRequestEcho(params.Logger),
		RecoverRequestEcho(params.Logger),
		middleware.Gzip(),
		// AuthenticateRequestEchoJWT[T](params.Config, params.Logger, params.KeyFuncJWT, params.PrincipalFactory),
	}
}

func LogRequestEcho(logger logging.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var logEvent logging.Event
			if v.Error != nil {
				logEvent = logger.Error()
			} else {
				logEvent = logger.Info()
			}
			bytesIn, _ := strconv.ParseInt(v.ContentLength, 10, 64)
			logEvent.
				WithField("id", v.RequestID).
				WithField("start_time", v.StartTime).
				WithField("remote_ip", v.RemoteIP).
				WithField("host", v.Host).
				WithField("method", v.Method).
				WithField("uri", v.URI).
				WithField("uri_path", v.URIPath).
				WithField("route_path", v.RoutePath).
				WithField("referer", v.Referer).
				WithField("user_agent", v.UserAgent).
				WithField("status", v.Status).
				WithField("error", v.Error).
				WithField("latency", v.Latency).
				WithField("latency_human", v.Latency.String()).
				WithField("protocol", v.Protocol).
				WithField("bytes_in", bytesIn).
				WithField("bytes_out", v.ResponseSize).
				Write("got request")
			return nil
		},
		HandleError:      true,
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
	})
}

func RecoverRequestEcho(logger logging.Logger) echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.WithError(err).WithField("stack", stack).Write("recovered from panic")
			return err
		},
	})
}

func AuthenticateRequestEchoJWT(cfg ConfigHTTP, logger logging.Logger, keyFuncJWT keyfunc.Keyfunc,
	factory security.PrincipalFactory[*jwt.Token]) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return cfg.AuthenticationWhitelistSet.Contains(c.Request().RequestURI)
		},
		SuccessHandler: func(c echo.Context) {
			// injects principal in context.Context. Using echo's context won't suffice
			// as PrincipalFactory relies on Go's context.
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				logger.
					WithError(errors.New("transport: cannot cast jwt token from echo context")).
					Write("could not create principal context")
				return
			}

			ctx, err := factory.NewContextWithPrincipal(c.Request().Context(), token)
			if err != nil {
				logger.WithError(err).Write("could not create principal context")
				return
			}
			req := c.Request().WithContext(ctx) // uses shallow copy, reducing extra malloc
			c.SetRequest(req)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return convertErrorEchoJWT(err)
		},
		SigningMethod: "RS256",
		KeyFunc:       keyFuncJWT.Keyfunc,
	})
}

// Authorization middlewares

func HasAnyAuthoritiesEcho(authorities ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := security.HasAnyAuthorities(c.Request().Context(), authorities); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func HasAuthoritiesEcho(authorities ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := security.HasAuthorities(c.Request().Context(), authorities); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func IsResourceOwnerEcho(resourceOwnerParamKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ownerID := c.Param(resourceOwnerParamKey)
			if err := security.IsResourceOwner(c.Request().Context(), ownerID); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func IsResourceOwnerOrHasAnyAuthoritiesEcho(resourceOwnerParamKey string,
	authorities ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ownerID := c.Param(resourceOwnerParamKey)
			if err := security.IsResourceOwnerOrHasAnyAuthorities(c.Request().Context(), ownerID, authorities); err != nil {
				return err
			}
			return next(c)
		}
	}
}
