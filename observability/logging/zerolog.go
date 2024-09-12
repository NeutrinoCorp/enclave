package logging

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog"

	"github.com/neutrinocorp/geck/application"
	"github.com/neutrinocorp/geck/observability/tracing"
)

// NewApplicationLogger allocates a zerolog.Logger instance with configuration.Application fields.
func NewApplicationLogger(cfg application.Config, writer io.Writer) zerolog.Logger {
	return zerolog.New(writer).With().
		Str("application_name", cfg.ApplicationName).
		Str("application_environment", cfg.Environment).
		Str("application_version", cfg.Version).
		Logger()
}

// ZerologEvent is zerolog the implementation of Event.
type ZerologEvent struct {
	ev *zerolog.Event
}

var _ Event = ZerologEvent{}

// WithField appends a field to the context.
func (z ZerologEvent) WithField(field string, val any) Event {
	switch v := val.(type) {
	case string:
		return ZerologEvent{ev: z.ev.Str(field, v)}
	case []string:
		return ZerologEvent{ev: z.ev.Strs(field, v)}
	case int:
		return ZerologEvent{ev: z.ev.Int(field, v)}
	case []int:
		return ZerologEvent{ev: z.ev.Ints(field, v)}
	case int8:
		return ZerologEvent{ev: z.ev.Int8(field, v)}
	case int16:
		return ZerologEvent{ev: z.ev.Int16(field, v)}
	case int32:
		return ZerologEvent{ev: z.ev.Int32(field, v)}
	case int64:
		return ZerologEvent{ev: z.ev.Int64(field, v)}
	case uint:
		return ZerologEvent{ev: z.ev.Uint(field, v)}
	case []uint:
		return ZerologEvent{ev: z.ev.Uints(field, v)}
	case uint8:
		return ZerologEvent{ev: z.ev.Uint8(field, v)}
	case uint16:
		return ZerologEvent{ev: z.ev.Uint16(field, v)}
	case uint32:
		return ZerologEvent{ev: z.ev.Uint32(field, v)}
	case uint64:
		return ZerologEvent{ev: z.ev.Uint64(field, v)}
	case bool:
		return ZerologEvent{ev: z.ev.Bool(field, v)}
	case float32:
		return ZerologEvent{ev: z.ev.Float32(field, v)}
	case []float32:
		return ZerologEvent{ev: z.ev.Floats32(field, v)}
	case float64:
		return ZerologEvent{ev: z.ev.Float64(field, v)}
	case []float64:
		return ZerologEvent{ev: z.ev.Floats64(field, v)}
	case time.Time:
		return ZerologEvent{ev: z.ev.Time(field, v)}
	case []time.Time:
		return ZerologEvent{ev: z.ev.Times(field, v)}
	case time.Duration:
		return ZerologEvent{ev: z.ev.Dur(field, v)}
	case []time.Duration:
		return ZerologEvent{ev: z.ev.Durs(field, v)}
	case []byte:
		return ZerologEvent{ev: z.ev.Bytes(field, v)}
	case net.IP:
		return ZerologEvent{ev: z.ev.IPAddr(field, v)}
	case net.IPNet:
		return ZerologEvent{ev: z.ev.IPPrefix(field, v)}
	case net.HardwareAddr:
		return ZerologEvent{ev: z.ev.MACAddr(field, v)}
	case error:
		return ZerologEvent{ev: z.ev.Err(v)}
	default:
		return ZerologEvent{ev: z.ev.Any(field, val)}
	}
}

// Write writes a new log entry into the Logger instance (most probably will write to an underlying io.Writer instance).
func (z ZerologEvent) Write(msg string) {
	z.ev.Msg(msg)
}

// WriteWithCtx writes a new log entry into the Logger instance (most probably will write to an underlying io.Writer instance).
//
// Uses context.Context to retrieve (and possibly append) useful information like trace identifiers.
func (z ZerologEvent) WriteWithCtx(ctx context.Context, msg string) {
	if traceID, _ := tracing.GetTraceIDFromContext(ctx); traceID != "" {
		z.WithField("trace_id", traceID)
	}
	z.Write(msg)
}

// ZerologLogger is the zerolog implementation of Logger.
type ZerologLogger struct {
	ModuleName string
	Logger     zerolog.Logger
}

var _ Logger = ZerologLogger{}

// NewZerologLogger allocates a new ZerologLogger instance.
func NewZerologLogger(l zerolog.Logger) ZerologLogger {
	return ZerologLogger{Logger: l}
}

// Level creates an Event context to write a new log entry.
func (z ZerologLogger) Level(lvl Level) Event {
	var event *zerolog.Event
	switch lvl {
	case DebugLevel:
		event = z.Logger.Debug()
	case InfoLevel:
		event = z.Logger.Info()
	case WarnLevel:
		event = z.Logger.Warn()
	case TraceLevel:
		event = z.Logger.Trace()
	case ErrorLevel:
		event = z.Logger.Error()
	default:
		event = z.Logger.Debug()
	}

	if z.ModuleName != "" {
		event = event.Str("module", z.ModuleName)
	}
	return ZerologEvent{ev: event}
}

// Module allocates a Logger instance with a module field.
func (z ZerologLogger) Module(name string) Logger {
	z.ModuleName = name
	return z
}

// Debug creates an Event context to write a new log entry with DebugLevel.
func (z ZerologLogger) Debug() Event {
	return z.Level(DebugLevel)
}

// Info creates an Event context to write a new log entry with InfoLevel.
func (z ZerologLogger) Info() Event {
	return z.Level(InfoLevel)
}

// Warn creates an Event context to write a new log entry with WarnLevel.
func (z ZerologLogger) Warn() Event {
	return z.Level(WarnLevel)
}

// Trace creates an Event context to write a new log entry with TraceLevel.
func (z ZerologLogger) Trace() Event {
	return z.Level(TraceLevel)
}

// Error creates an Event context to write a new log entry with ErrorLevel.
func (z ZerologLogger) Error() Event {
	return z.Level(ErrorLevel)
}

// WithError creates an Event context to write a new log entry with ErrorLevel and appends an `error` field.
func (z ZerologLogger) WithError(err error) Event {
	return z.Level(ErrorLevel).WithField("error", err)
}
