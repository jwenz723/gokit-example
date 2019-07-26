package transport

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/addsvc/pkg/addservice"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/ratelimit"
	"github.com/jwenz723/kit-mw/eplogger"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	SumEndpoint    endpoint.Endpoint
	ConcatEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc addservice.Service, logger log.Logger) Set {
	var sumEndpoint endpoint.Endpoint
	{
		l := log.With(logger, "method", "Sum")
		sumEndpoint = MakeSumEndpoint(svc)
		sumEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(sumEndpoint)
		sumEndpoint = eplogger.LoggingMiddleware(level.Info(l), level.Error(l))(sumEndpoint)
	}
	var concatEndpoint endpoint.Endpoint
	{
		l := log.With(logger, "method", "Concat")
		concatEndpoint = MakeConcatEndpoint(svc)
		concatEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(concatEndpoint)
		concatEndpoint = eplogger.LoggingMiddleware(level.Info(l), level.Error(l))(concatEndpoint)
	}
	return Set{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

// Sum implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Sum(ctx context.Context, a, b int) (int, error) {
	resp, err := s.SumEndpoint(ctx, SumRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(SumResponse)
	return response.V, response.Err
}

// Concat implements the service interface, so Set may be used as a
// service. This is primarily useful in the context of a client library.
func (s Set) Concat(ctx context.Context, a, b string) (string, error) {
	resp, err := s.ConcatEndpoint(ctx, ConcatRequest{A: a, B: b})
	if err != nil {
		return "", err
	}
	response := resp.(ConcatResponse)
	return response.V, response.Err
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeSumEndpoint(s addservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		return SumResponse{V: v, Err: err}, nil
	}
}

// MakeConcatEndpoint constructs a Concat endpoint wrapping the service.
func MakeConcatEndpoint(s addservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ConcatRequest)
		v, err := s.Concat(ctx, req.A, req.B)
		return ConcatResponse{V: v, Err: err}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = SumResponse{}
	_ endpoint.Failer = ConcatResponse{}

	// eplogger.AppendKeyvalser interface
	_ eplogger.AppendKeyvalser = ConcatRequest{}
	_ eplogger.AppendKeyvalser = ConcatResponse{}
	_ eplogger.AppendKeyvalser = SumRequest{}
	_ eplogger.AppendKeyvalser = SumResponse{}
)

// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A, B int
}

// AppendKeyvals implements AppendKeyvalser to return keyvals specific to SumRequest for logging
func (s SumRequest) AppendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals,
		"SumRequest.A", s.A,
		"SumRequest.B", s.B)
}

// SumResponse collects the response values for the Sum method.
type SumResponse struct {
	V   int   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// AppendKeyvals implements AppendKeyvalser to return keyvals specific to SumResponse for logging
func (s SumResponse) AppendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals,
		"SumResponse.R", s.V,
		"SumResponse.Err", s.Err)
}

// Failed implements endpoint.Failer.
func (r SumResponse) Failed() error { return r.Err }

// ConcatRequest collects the request parameters for the Concat method.
type ConcatRequest struct {
	A, B string
}

// AppendKeyvals implements AppendKeyvalser to return keyvals specific to ConcatRequest for logging
func (c ConcatRequest) AppendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals,
		"ConcatRequest.A", c.A,
		"ConcatRequest.B", c.B)
}

// ConcatResponse collects the response values for the Concat method.
type ConcatResponse struct {
	V   string `json:"v"`
	Err error  `json:"-"`
}

// AppendKeyvals implements AppendKeyvalser to return keyvals specific to ConcatResponse for logging
func (c ConcatResponse) AppendKeyvals(keyvals []interface{}) []interface{} {
	return append(keyvals,
		"ConcatResponse.V", c.V,
		"ConcatResponse.Err", c.Err)
}

// Failed implements endpoint.Failer.
func (r ConcatResponse) Failed() error { return r.Err }
