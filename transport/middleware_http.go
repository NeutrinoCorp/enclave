package transport

import (
	"context"
	"errors"
	"strconv"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/neutrinocorp/geck/observability/logging"
	"github.com/neutrinocorp/geck/observability/tracing"
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

type TraceIDEchoParams struct {
	fx.In
	TraceFactory tracing.TraceFactory `optional:"true"`
}

// NewTraceIDEcho appends a trace identifier to each request using tracing.NewTracedContext.
func NewTraceIDEcho(params TraceIDEchoParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// injects principal in context.Context. Using echo's context won't suffice
			// as GECK tracing package relies on Go's context.
			var ctx context.Context
			if params.TraceFactory != nil {
				ctx = params.TraceFactory.NewTracedContext(c.Request().Context())
			} else {
				ctx = tracing.NewTracedContext(c.Request().Context())
			}
			req := c.Request().WithContext(ctx) // uses shallow copy, reducing extra malloc
			c.SetRequest(req)
			return next(c)
		}
	}
}

type DefaultEchoMiddlewareParams struct {
	fx.In

	Config ConfigHTTP
	Logger logging.Logger
}

// NewDefaultEchoMiddlewareGroup allocates an Echo middleware group. An array is returned to guarantee
// ordering within this group. Any other middlewares outside this group will be injected into application in
// a non-deterministic way (regarding ordering).
func NewDefaultEchoMiddlewareGroup(params DefaultEchoMiddlewareParams, paramsTrace TraceIDEchoParams) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		NewTraceIDEcho(paramsTrace),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			TargetHeader: params.Config.RequestIDTargetHeader,
		}),
		NewLogRequestEcho(params.Logger),
		NewRecoverRequestEcho(params.Logger),
		middleware.Gzip(),
	}
}

func NewLogRequestEcho(logger logging.Logger) echo.MiddlewareFunc {
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
				WithField("request_id", v.RequestID).
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
				WriteWithCtx(c.Request().Context(), "got request")
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

func NewRecoverRequestEcho(logger logging.Logger) echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisableStackAll: true,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.WithError(err).WithField("stack", stack).WriteWithCtx(c.Request().Context(), "recovered from panic")
			return err
		},
	})
}

type NewEchoJWTAuthenticatorConfigParams struct {
	fx.In

	Config           security.ConfigJWT
	ServerConfig     ConfigHTTP
	Logger           logging.Logger
	PrincipalFactory security.PrincipalFactory[*jwt.Token]
	KeyFunc          keyfunc.Keyfunc `optional:"true"`
}

func NewEchoJWTAuthenticatorConfig(params NewEchoJWTAuthenticatorConfigParams) echojwt.Config {
	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return params.ServerConfig.AuthenticationWhitelistSet.Contains(c.Request().RequestURI)
		},
		SuccessHandler: func(c echo.Context) {
			// injects principal in context.Context. Using echo's context won't suffice
			// as PrincipalFactory relies on Go's context.
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				params.Logger.
					WithError(errors.New("transport: cannot cast jwt token from echo context")).
					WriteWithCtx(c.Request().Context(), "could not create principal context")
				return
			}

			ctx, err := params.PrincipalFactory.NewContextWithPrincipal(c.Request().Context(), token)
			if err != nil {
				params.Logger.WithError(err).WriteWithCtx(c.Request().Context(), "could not create principal context")
				return
			}
			req := c.Request().WithContext(ctx) // uses shallow copy, reducing extra malloc
			c.SetRequest(req)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return convertErrorEchoJWT(c.Request().Context(), params.Logger, err)
		},
		SigningMethod: params.Config.SigningMethod,
		SigningKey:    params.Config.SigningKey,
		SigningKeys: lo.MapEntries(params.Config.SigningKeys, func(key string, value string) (string, any) {
			return key, value
		}),
		KeyFunc: params.KeyFunc.Keyfunc,
	}
}

func NewEchoJWTAuthenticator(config echojwt.Config) echo.MiddlewareFunc {
	return echojwt.WithConfig(config)
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
