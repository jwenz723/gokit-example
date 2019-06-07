package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/jwenz723/logwrapper/pkg/service"
)

// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// LoggingKeyvals implements LoggingKeyvalser to return keyvals specific to SumRequest for logging
func (s SumRequest) LoggingKeyvals() []interface{} {
	return []interface{}{"SumRequest.A", s.A, "SumRequest.B", s.B}
}

// SumResponse collects the response parameters for the Sum method.
type SumResponse struct {
	R   int   `json:"r"`
	Err error `json:"err"`
}

// LoggingKeyvals implements LoggingKeyvalser to return keyvals specific to SumResponse for logging
func (s SumResponse) LoggingKeyvals() []interface{} {
	return []interface{}{"SumResponse.R", s.R, "SumResponse.Err", s.Err}
}

// MakeSumEndpoint returns an endpoint that invokes Sum on the service.
func MakeSumEndpoint(s service.LogwrapperService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		r, err := s.Sum(ctx, req.A, req.B)
		return SumResponse{
			Err: err,
			R:   r,
		}, nil
	}
}

// Failed implements Failer.
func (r SumResponse) Failed() error {
	return r.Err
}

// MultiplyRequest collects the request parameters for the Multiply method.
type MultiplyRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// LoggingKeyvals implements LoggingKeyvalser to return keyvals specific to MultiplyRequest for logging
func (s MultiplyRequest) LoggingKeyvals() []interface{} {
	return []interface{}{"MultiplyRequest.A", s.A, "MultiplyRequest.B", s.B}
}

// MultiplyResponse collects the response parameters for the Multiply method.
type MultiplyResponse struct {
	R   int   `json:"r"`
	Err error `json:"err"`
}

// LoggingKeyvals implements LoggingKeyvalser to return keyvals specific to MultiplyRequest for logging
func (s MultiplyResponse) LoggingKeyvals() []interface{} {
	return []interface{}{"MultiplyResponse.R", s.R, "MultiplyResponse.Err", s.Err}
}

// MakeMultiplyEndpoint returns an endpoint that invokes Multiply on the service.
func MakeMultiplyEndpoint(s service.LogwrapperService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MultiplyRequest)
		r, err := s.Multiply(ctx, req.A, req.B)
		return MultiplyResponse{
			Err: err,
			R:   r,
		}, nil
	}
}

// Failed implements Failer.
func (r MultiplyResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Sum implements Service. Primarily useful in a client.
func (e Endpoints) Sum(ctx context.Context, a int, b int) (r int, err error) {
	request := SumRequest{
		A: a,
		B: b,
	}
	response, err := e.SumEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SumResponse).R, response.(SumResponse).Err
}

// Multiply implements Service. Primarily useful in a client.
func (e Endpoints) Multiply(ctx context.Context, a int, b int) (r int, err error) {
	request := MultiplyRequest{
		A: a,
		B: b,
	}
	response, err := e.MultiplyEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(MultiplyResponse).R, response.(MultiplyResponse).Err
}
