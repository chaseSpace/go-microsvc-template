package xgrpc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microsvc/pkg/xerr"
	"microsvc/pkg/xlog"
	"strings"
	"time"
)

// -------------------------- clientUtil ------------------------------

type clientUtil struct{}

var cutil = clientUtil{}

func (clientUtil) newCircuitBreaker(name string) *gobreaker.CircuitBreaker {
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
func (clientUtil) fromGatewayCall(ctx context.Context) bool {
	return GetOutgoingMdVal(ctx, MdKeyFromGatewayFlag) == MdKeyFlagExist
}

func (clientUtil) beautifyReqAndResInClient(ctx context.Context, req, reply interface{}) (interface{}, interface{}) {
	if !cutil.fromGatewayCall(ctx) {
		return req, reply
	}
	// When grpc call is from gateway,the request type must be `[]byte`, and the reply type must be `*bytes.Buffer`
	req = string(req.([]byte))
	reply = reply.(*bytes.Buffer).String()
	return req, reply
}

func (clientUtil) breakerTakeError(err error) bool {
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

// -------------------------- serverUtil ------------------------------
type serverUtil struct{}

var sutil = serverUtil{}

func (serverUtil) setupCtx(ctx context.Context, method string) (context.Context, error) {
	isExtMethod := false
	fromGateway := GetIncomingMdVal(ctx, MdKeyFromGatewayFlag) == MdKeyFlagExist

	// method such as: /svc.user.UserExt/Signup
	ss := strings.Split(method, "/")
	if len(ss) == 3 {
		if strings.HasSuffix(ss[1], "Ext") {
			isExtMethod = true
		} else if !strings.HasSuffix(ss[1], "Int") {
			return nil, fmt.Errorf("illegal grpc method: %s", method)
		}
		ctx = context.WithValue(ctx, CtxServerSideKey{}, CtxServerSideVal{
			IsExtMethod: isExtMethod,
			FromGateway: fromGateway,
		})
		return ctx, nil
	}
	return nil, fmt.Errorf("illegal grpc method: %s", method)
}

func (serverUtil) isExtMethod(ctx context.Context) bool {
	val := ctx.Value(CtxServerSideKey{}).(CtxServerSideVal)
	return val.IsExtMethod
}

func (serverUtil) fromGatewayCall(ctx context.Context) bool {
	return GetIncomingMdVal(ctx, MdKeyFromGatewayFlag) == MdKeyFlagExist
}
