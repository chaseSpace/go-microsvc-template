package xerr

// Most common errors are defined referenced by HTTP status codes.

// 200-400 series
var (
	ErrOK          = XErr{Code: 200, Msg: "OK"}
	ErrParams      = XErr{Code: 400, Msg: "ErrParams"}
	ErrApiNotFound = XErr{Code: 404, Msg: "ErrApiNotFound"}

	ErrBadRequest = ErrParams.NewMsg("ErrBadRequest")
)

// 500 series
var (
	ErrInternal       = XErr{Code: 500, Msg: "ErrInternal"}
	ErrGatewayTimeout = XErr{Code: 504, Msg: "ErrGatewayTimeout"}
	ErrNoRPCClient    = ErrInternal.NewMsg("no available rpc client")
)

// Customized errors
var (
	_          = 0
	ErrTimeout = XErr{Code: 1000, Msg: "ErrTimeout"}
)
