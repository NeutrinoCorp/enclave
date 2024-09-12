package tracing

import (
	"context"

	"github.com/neutrinocorp/geck/identifier"
)

var (
	tracerFactoryID identifier.FactoryUUID
)

type TraceIDContextType string

const TraceIDContextKey TraceIDContextType = "tracing.trace_id"

// NewTracedContext appends a generated id to the given context.Context.
//
// Uses identifier.FactoryUUID as identifier generation algorithm.
func NewTracedContext(ctx context.Context) context.Context {
	id, _ := tracerFactoryID.NewIdentifier()
	return context.WithValue(ctx, TraceIDContextKey, id)
}

// GetTraceIDFromContext retrieves a trace identifier from the given context.Context. Produces ErrTraceIDNotFound if
// not found or id cannot be cast.
func GetTraceIDFromContext(ctx context.Context) (string, error) {
	traceID, ok := ctx.Value(TraceIDContextKey).(string)
	if !ok {
		return "", ErrTraceIDNotFound
	}
	return traceID, nil
}
