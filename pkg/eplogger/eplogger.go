package eplogger

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"time"
)

type LogKeyvalsAdder interface {
	AddLogKeyvals(log.Logger) log.Logger
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, the resulting error (if any), and
// keyvals specific to the request and response object if they implement
// the LogKeyvalsAdder interface.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				// Check if request implements the LogKeyvalsAdder interface
				if l, ok := request.(LogKeyvalsAdder); ok {
					// Update logger to contain keyvals specific to request
					logger = l.AddLogKeyvals(logger)
				}
				// Check if response implements the LogKeyvalsAdder interface
				if l, ok := response.(LogKeyvalsAdder); ok {
					// Update logger to contain keyvals specific to request
					logger = l.AddLogKeyvals(logger)
				}
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}
