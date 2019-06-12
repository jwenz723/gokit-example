package eplogger

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type LoggingKeyvalser interface {
	LoggingKeyvals() (keyvals []interface{})
}

const (
	tookKey     = "took"
	transErrKey = "transport_error"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, the resulting error (if any), and
// keyvals specific to the request and response object if they implement
// the LoggingKeyvalser interface.
//
// The level specified as defaultLevel will be used when the resulting error
// is nil otherwise level.Error will be used.
func LoggingMiddleware(logger log.Logger, defaultLevel level.Value) endpoint.Middleware {
	// This will set a default log level if one is not set when logger.Log() is executed
	logger = level.NewInjector(logger, defaultLevel)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				if err != nil {
					logger = level.Error(logger)
				}
				logger = logWithLoggingKeyvalser(logger, request)
				logger = logWithLoggingKeyvalser(logger, response)
				logger = logWithError(logger, err)
				logger = logWithDuration(logger, time.Since(begin))
				logger.Log()
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// logWithLoggingKeyvalser will add the keyvals returned by k.LoggingKeyvals
// into logger if k implements the LoggingKeyvalser interface
func logWithLoggingKeyvalser(logger log.Logger, k interface{}) log.Logger {
	if l, ok := k.(LoggingKeyvalser); ok {
		return log.With(logger, l.LoggingKeyvals()...)
	}
	return logger
}

// logWithError will add err into the keyvals of logger
func logWithError(logger log.Logger, err error) log.Logger {
	return log.With(logger, transErrKey, err)
}

// logWithDuration will add d into the keyvals of logger
func logWithDuration(logger log.Logger, d time.Duration) log.Logger {
	return log.With(logger, tookKey, d)
}
