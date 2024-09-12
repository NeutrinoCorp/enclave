package logging

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/neutrinocorp/geck/observability/tracing"
)

// StdEvent is standard library (log.Logger) the implementation of Event.
type StdEvent struct {
	level  Level
	logger *log.Logger
	module string
	fields map[string]any
}

var _ Event = &StdEvent{}

// WithField appends a field to the context.
func (s *StdEvent) WithField(field string, val any) Event {
	s.fields[field] = val
	return s
}

// Write writes a new log entry into the Logger instance (most probably will write to an underlying io.Writer instance).
func (s *StdEvent) Write(msg string) {
	buf := strings.Builder{}
	i := 0
	for k, v := range s.fields {
		switch val := v.(type) {
		case string:
			buf.WriteString(fmt.Sprintf("%s:%q", k, val))
		default:
			buf.WriteString(fmt.Sprintf("%s:%v", k, v))
		}

		if i < len(s.fields)-1 {
			buf.WriteByte(' ')
		}
		i++
	}

	var lvl string
	switch s.level {
	case DebugLevel:
		lvl = "DEBUG"
	case InfoLevel:
		lvl = "INFO"
	case WarnLevel:
		lvl = "WARN"
	case TraceLevel:
		lvl = "TRACE"
	case ErrorLevel:
		lvl = "ERROR"
	default:
		lvl = "DEBUG"
	}

	s.logger.Printf("%s %s message:%q", lvl, buf.String(), msg)
}

// WriteWithCtx writes a new log entry into the Logger instance (most probably will write to an underlying io.Writer instance).
//
// Uses context.Context to retrieve (and possibly append) useful information like trace identifiers.
func (s *StdEvent) WriteWithCtx(ctx context.Context, msg string) {
	if traceID, _ := tracing.GetTraceIDFromContext(ctx); traceID != "" {
		s.WithField("trace_id", traceID)
	}
	s.Write(msg)
}

// StdLogger is the standard library (log.Logger) implementation of Logger.
type StdLogger struct {
	ModuleName string
	Logger     *log.Logger
}

var _ Logger = StdLogger{}

// NewStdLogger allocates a new StdLogger instance.
func NewStdLogger(l *log.Logger) StdLogger {
	return StdLogger{Logger: l}
}

// Level creates an Event context to write a new log entry.
func (s StdLogger) Level(lvl Level) Event {
	ev := &StdEvent{
		level:  lvl,
		logger: s.Logger,
		module: s.ModuleName,
		fields: map[string]any{},
	}
	if s.ModuleName != "" {
		ev.WithField("module", s.ModuleName)
	}
	return ev
}

// Module allocates a Logger instance with a module field.
func (s StdLogger) Module(name string) Logger {
	s.ModuleName = name
	return s
}

// Debug creates an Event context to write a new log entry with DebugLevel.
func (s StdLogger) Debug() Event {
	return s.Level(DebugLevel)
}

// Info creates an Event context to write a new log entry with InfoLevel.
func (s StdLogger) Info() Event {
	return s.Level(InfoLevel)
}

// Warn creates an Event context to write a new log entry with WarnLevel.
func (s StdLogger) Warn() Event {
	return s.Level(WarnLevel)
}

// Trace creates an Event context to write a new log entry with TraceLevel.
func (s StdLogger) Trace() Event {
	return s.Level(TraceLevel)
}

// Error creates an Event context to write a new log entry with ErrorLevel.
func (s StdLogger) Error() Event {
	return s.Level(ErrorLevel)
}

// WithError creates an Event context to write a new log entry with ErrorLevel and appends an `error` field.
func (s StdLogger) WithError(err error) Event {
	return s.Error().WithField("error", err.Error())
}
