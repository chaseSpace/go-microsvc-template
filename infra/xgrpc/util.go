package xgrpc

import (
	"bytes"
	"fmt"
	"github.com/samber/lo"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"time"
)

func beautifyReqAndResInClient(req, reply interface{}) (interface{}, interface{}) {
	// When grpc call is from gateway,the request is []byte, and the reply is *bytes.Buffer
	_req, _ := req.([]byte)
	if _req != nil {
		req = string(_req)
	}
	res, _ := reply.(*bytes.Buffer)
	if res != nil {
		reply = res.String()
	}
	return req, reply
}

func newCircuitBreaker(name string) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        name,
		MaxRequests: 1,                //  maximum number of requests allowed to pass through when the breaker is half-open
		Interval:    time.Second * 30, // cyclic period of breaker to clear interval counter, defaults to 0 that indicates the breaker never clear interval counter
		Timeout:     time.Second * 10, // timeout for CircuitBreaker stay open, breaker switch to half-open after `timeout`, default 60s.
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// define the condition of the breaker gets to open state
			//fmt.Printf("ReadyToTrip_xx  %+v\n", counts)
			return counts.ConsecutiveFailures > 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			xlog.Warn(fmt.Sprintf("grpc call - circuit breaker state change: %s, %s -> %s", name, from, to))
		},
	})
}

func breakerTakeError(err error) bool {
	if xe, ok := err.(xerr.XErr); ok && xe.IsInternal() {
		return true
	}
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return lo.Contains([]codes.Code{
		codes.Unavailable,
		codes.DeadlineExceeded,
		codes.Aborted,
		codes.FailedPrecondition,
	}, s.Code())
}
