package service

import (
	"context"

	log "github.com/go-kit/kit/log"
)

type Middleware func(LogwrapperService) LogwrapperService

type loggingMiddleware struct {
	logger log.Logger
	next   LogwrapperService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next LogwrapperService) LogwrapperService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Sum(ctx context.Context, a int, b int) (r int, err error) {
	defer func() {
		l.logger.Log("method", "Sum", "a", a, "b", b, "r", r, "err", err)
	}()
	return l.next.Sum(ctx, a, b)
}
func (l loggingMiddleware) Multiply(ctx context.Context, a int, b int) (r int, err error) {
	defer func() {
		l.logger.Log("method", "Multiply", "a", a, "b", b, "r", r, "err", err)
	}()
	return l.next.Multiply(ctx, a, b)
}
