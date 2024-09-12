package tracing

import "errors"

// ErrTraceIDNotFound no trace id was found.
var ErrTraceIDNotFound = errors.New("trace id not found")
