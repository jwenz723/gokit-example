package eplogger

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
)

const stringFieldKey = "StringField"

type LoggingKeyvalserTest struct {
	StringField string
}

func (l LoggingKeyvalserTest) LoggingKeyvals() (keyvals []interface{}) {
	return []interface{}{stringFieldKey, l.StringField}
}

func TestLoggingMiddleware(t *testing.T) {
	var output []interface{}

	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = keyvals
		return nil
	}))

	req := LoggingKeyvalserTest{
		StringField: "req string",
	}

	resp := LoggingKeyvalserTest{
		StringField: "resp string",
	}

	nextExecuted := false

	nextFunc := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		nextExecuted = true
		return resp, nil
	}
	ep := LoggingMiddleware(logger)(nextFunc)
	ep(context.Background(), req)

	for i, v := range output {
		if v == req.StringField && output[i-1] != stringFieldKey {
			t.Errorf("invalid req StringField value")
		}

		if v == resp.StringField && output[i-1] != stringFieldKey {
			t.Errorf("invalid resp StringField value")
		}

		if v == transErrKey && output[i+1] != nil {
			t.Errorf("unexpected value for %s", transErrKey)
		}

		if !nextExecuted {
			t.Errorf("nextFunc was never executed")
		}
	}
}
