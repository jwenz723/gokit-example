package service

import "context"

// LogwrapperService describes the service.
type LogwrapperService interface {
	Sum(ctx context.Context, a, b int) (r int, err error)
	Multiply(ctx context.Context, a, b int) (r int, err error)
}

type basicLogwrapperService struct{}

func (ba *basicLogwrapperService) Sum(ctx context.Context, a int, b int) (r int, err error) {
	return a + b, nil
}
func (ba *basicLogwrapperService) Multiply(ctx context.Context, a int, b int) (r int, err error) {
	return a * b, nil
}

// NewBasicLogwrapperService returns a naive, stateless implementation of LogwrapperService.
func NewBasicLogwrapperService() LogwrapperService {
	return &basicLogwrapperService{}
}

// New returns a LogwrapperService with all of the expected middleware wired in.
func New(middleware []Middleware) LogwrapperService {
	var svc LogwrapperService = NewBasicLogwrapperService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
