package eplogger

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"testing"
	"time"
)

const stringFieldKey = "StringField"

// LoggingKeyvalserTest is a type that implements LoggingKeyvalser interface used for testing
type LoggingKeyvalserTest struct {
	StringField string
}

// LoggingKeyvals implements LoggingKeyvalser interface
func (l LoggingKeyvalserTest) LoggingKeyvals() (keyvals []interface{}) {
	return []interface{}{stringFieldKey, l.StringField}
}

func TestLoggingMiddleware(t *testing.T) {
	var tests = map[string]struct {
		expectLevel          level.Value
		inReq                string
		inResp               string
		inRespErr            error
		withLoggingKeyvalser bool
	}{
		"nil error": {
			expectLevel:          level.InfoValue(),
			inReq:                "req string",
			inResp:               "resp string",
			inRespErr:            nil,
			withLoggingKeyvalser: true,
		},
		"non-nil error": {
			expectLevel:          level.ErrorValue(),
			inReq:                "req string",
			inResp:               "resp string",
			inRespErr:            errors.New("an error"),
			withLoggingKeyvalser: true,
		},
		"nil error, no LoggingKeyvalser": {
			expectLevel:          level.InfoValue(),
			inReq:                "req string",
			inResp:               "resp string",
			inRespErr:            nil,
			withLoggingKeyvalser: false,
		},
		"non-nil error, no LoggingKeyvalser": {
			expectLevel:          level.ErrorValue(),
			inReq:                "req string",
			inResp:               "resp string",
			inRespErr:            errors.New("an error"),
			withLoggingKeyvalser: false,
		},
	}

	var output []interface{}
	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = keyvals
		return nil
	}))

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			var req, resp interface{}
			if tt.withLoggingKeyvalser {
				req = LoggingKeyvalserTest{
					StringField: tt.inReq,
				}
				resp = LoggingKeyvalserTest{
					StringField: tt.inResp,
				}
			} else {
				req = tt.inReq
				resp = tt.inResp
			}

			// Simulate a go-kit endpoint
			endpointExecuted := false
			ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
				endpointExecuted = true
				return resp, tt.inRespErr
			}

			// Wrap the simulated endpoint with the middleware
			epWithMw := LoggingMiddleware(logger, level.InfoValue())(ep)

			// Execute the endpoint and middleware
			epWithMw(context.Background(), req)

			if tt.withLoggingKeyvalser {
				if want, have := "level", output[0]; want != have {
					t.Errorf("output[0]: want %s, have %s", want, have)
				}
				if want, have := tt.expectLevel, output[1]; want != have {
					t.Errorf("output[1]: want %s, have %s", want, have)
				}
				if want, have := stringFieldKey, output[2]; want != have {
					t.Errorf("output[2]: want %s, have %s", want, have)
				}
				if want, have := tt.inReq, output[3]; want != have {
					t.Errorf("output[3]: want %s, have %s", want, have)
				}
				if want, have := stringFieldKey, output[4]; want != have {
					t.Errorf("output[4]: want %s, have %s", want, have)
				}
				if want, have := tt.inResp, output[5]; want != have {
					t.Errorf("output[5]: want %s, have %s", want, have)
				}
				if want, have := transErrKey, output[6]; want != have {
					t.Errorf("output[6]: want %s, have %s", want, have)
				}
				if want, have := tt.inRespErr, output[7]; want != have {
					t.Errorf("output[7]: want %s, have %s", want, have)
				}
				if want, have := tookKey, output[8]; want != have {
					t.Errorf("output[8]: want %s, have %s", want, have)
				}
				_, ok := output[9].(time.Duration)
				if !ok {
					t.Fatalf("want time.Time, have %T", output[9])
				}
				if !endpointExecuted {
					t.Errorf("endpoint was never executed")
				}
			} else {
				if want, have := "level", output[0]; want != have {
					t.Errorf("output[0]: want %s, have %s", want, have)
				}
				if want, have := tt.expectLevel, output[1]; want != have {
					t.Errorf("output[1]: want %s, have %s", want, have)
				}
				if want, have := transErrKey, output[2]; want != have {
					t.Errorf("output[2]: want %s, have %s", want, have)
				}
				if want, have := tt.inRespErr, output[3]; want != have {
					t.Errorf("output[3]: want %s, have %s", want, have)
				}
				if want, have := tookKey, output[4]; want != have {
					t.Errorf("output[4]: want %s, have %s", want, have)
				}
				_, ok := output[5].(time.Duration)
				if !ok {
					t.Fatalf("want time.Time, have %T", output[5])
				}
				if !endpointExecuted {
					t.Errorf("endpoint was never executed")
				}
			}
		})
	}
}

func TestLogWithLoggingKeyvalser(t *testing.T) {
	var output []interface{}
	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = keyvals
		return nil
	}))

	// Test that a type that implements LoggingKeyvalser interface will have values logged
	l := LoggingKeyvalserTest{
		StringField: "inReq string",
	}
	keysLogger := logWithLoggingKeyvalser(logger, l)
	keysLogger.Log()
	s, ok := output[1].(string)
	if !ok {
		t.Fatalf("want string, have %T", output[1])
	}
	if want, have := l.StringField, s; want != have {
		t.Errorf("output[1]: want %v, have %v", want, have)
	}

	// Test that a type that does NOT implement LoggingKeyvalser interface will not be logged
	m := "my value that shouldn't be logged"
	emptyLogger := logWithLoggingKeyvalser(logger, m)
	emptyLogger.Log()
	if len(output) > 0 {
		t.Errorf("output should be empty, have %v", output)
	}
}

func TestLogWithError(t *testing.T) {
	var output []interface{}
	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = keyvals
		return nil
	}))

	expectedErr := errors.New("an error")
	logger = logWithError(logger, expectedErr)
	logger.Log()

	err, ok := output[1].(error)
	if !ok {
		t.Fatalf("want error, have %T", output[1])
	}
	if want, have := expectedErr, err; want != have {
		t.Errorf("output[1]: want %v, have %v", want, have)
	}
}

func TestLogWithDuration(t *testing.T) {
	var output []interface{}
	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = keyvals
		return nil
	}))

	expectedDuration := time.Since(time.Now())
	logger = logWithDuration(logger, expectedDuration)
	logger.Log()

	duration, ok := output[1].(time.Duration)
	if !ok {
		t.Fatalf("want time.Duration, have %T", output[1])
	}
	if want, have := expectedDuration, duration; want != have {
		t.Errorf("output[1]: want %v, have %v", want, have)
	}
}

// BenchmarkLoggingMiddleware tests how long the middleware takes to execute when
// the resulting err is nil.
// The benchmark output by BenchmarkLoggingMiddlewareCreation should be subtracted from the
// benchmark output by this func.
func BenchmarkLoggingMiddleware(b *testing.B) {
	req := LoggingKeyvalserTest{
		StringField: "inReq string",
	}
	resp := LoggingKeyvalserTest{
		StringField: "inResp string",
	}

	// Simulate a go-kit endpoint
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return resp, nil
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Wrap the simulated endpoint with the middleware. Need to do this for each
		// b.N iteration so that a new logger instance can be passed to the middleware
		// to avoid memory copying from slowing the benchmark down during log.With()
		epWithMw := LoggingMiddleware(log.NewNopLogger(), level.InfoValue())(ep)

		// Execute the endpoint and middleware
		epWithMw(context.Background(), req)
	}
}

// BenchmarkLoggingMiddlewareWithErr tests how long the middleware takes to execute when
// the resulting err is not nil.
// The benchmark output by BenchmarkLoggingMiddlewareCreation should be subtracted from the
// benchmark output by this func.
func BenchmarkLoggingMiddlewareWithErr(b *testing.B) {
	req := LoggingKeyvalserTest{
		StringField: "inReq string",
	}
	resp := LoggingKeyvalserTest{
		StringField: "inResp string",
	}

	// Simulate a go-kit endpoint
	epErr := errors.New("an error")
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return resp, epErr
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Wrap the simulated endpoint with the middleware. Need to do this for each
		// b.N iteration so that a new logger instance can be passed to the middleware
		// to avoid memory copying from slowing the benchmark down during log.With()
		epWithMw := LoggingMiddleware(log.NewNopLogger(), level.InfoValue())(ep)

		// Execute the endpoint and middleware
		epWithMw(context.Background(), req)
	}
}

// BenchmarkLoggingMiddlewareCreation outputs how long initialization takes
// of the LoggingMiddleware. The benchmark returned here can be subtracted from
// other benchmarks to get an accurate representation of their purposes.
func BenchmarkLoggingMiddlewareCreation(b *testing.B) {
	resp := LoggingKeyvalserTest{
		StringField: "inResp string",
	}

	// Simulate a go-kit endpoint
	epErr := errors.New("an error")
	ep := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return resp, epErr
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Wrap the simulated endpoint with the middleware. Need to do this for each
		// b.N iteration so that a new logger instance can be passed to the middleware
		// to avoid memory copying from slowing the benchmark down during log.With()
		LoggingMiddleware(log.NewNopLogger(), level.InfoValue())(ep)
	}
}
