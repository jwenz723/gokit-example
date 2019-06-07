package endpoint

import (
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"testing"
)

func BenchmarkSumRequest_LoggingKeyvals(b *testing.B) {
	b.ReportAllocs()
	logger := log.NewLogfmtLogger(ioutil.Discard)
	r := SumRequest{
		A: 3,
		B: 4,
	}
	for i := 0; i < b.N; i++ {
		l := log.With(logger, r.LoggingKeyvals()...)
		l.Log()
	}
}
