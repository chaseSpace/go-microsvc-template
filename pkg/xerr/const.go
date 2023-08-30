package xerr

// Common error
var (
	_              = 0
	ErrOK          = XErr{Code: 200, Msg: "OK"}
	ErrParams      = XErr{Code: 400, Msg: "ErrParams"}
	ErrApiNotFound = XErr{Code: 404, Msg: "ErrApiNotFound"}
	ErrInternal    = XErr{Code: 500, Msg: "ErrInternal"}
	ErrTimeout     = XErr{Code: 1000, Msg: "ErrTimeout"}
)

// Internal error
var (
	_              = 0
	ErrNoRPCClient = ErrInternal.NewMsg("no available rpc client")
)
