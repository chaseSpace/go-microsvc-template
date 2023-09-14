package xerr

// Most common errors are defined referenced by HTTP status codes.

// 200-400 series
var (
	ErrNil         = XErr{Code: 200, Msg: "OK"}
	ErrParams      = XErr{Code: 400, Msg: "ErrParams"}
	ErrApiNotFound = XErr{Code: 404, Msg: "ErrApiNotFound"}

	ErrBadRequest = ErrParams.NewMsg("ErrBadRequest")
)

// 500 series
var (
	ErrInternal    = XErr{Code: 500, Msg: "ErrInternal"}
	ErrGRPCTimeout = XErr{Code: 504, Msg: "ErrGRPCTimeout"}
	ErrNoRPCClient = ErrInternal.NewMsg("no available rpc client")
)

// Customized errors
var (
	_             = 0
	ErrBizTimeout = XErr{Code: 1000, Msg: "ErrBizTimeout"}
	ErrThirdParty = XErr{Code: 1001, Msg: "ErrThirdParty"}
)
