package xerr

// Most common errors are defined referenced by HTTP status codes.

// 200-400 series
var (
	ErrNil              = XErr{Code: 200, Msg: "OK"}
	ErrParams           = XErr{Code: 400, Msg: "ErrParams"}
	ErrUnauthorized     = XErr{Code: 401, Msg: "ErrUnauthorized"}
	ErrForbidden        = XErr{Code: 403, Msg: "ErrForbidden"}
	ErrNotFound         = XErr{Code: 404, Msg: "ErrNotFound"}
	ErrMethodNotAllowed = XErr{Code: 405, Msg: "ErrMethodNotAllowed"}
	ErrReqTimeout       = XErr{Code: 408, Msg: "ErrReqTimeout"}

	ErrBadRequest = ErrParams.NewMsg("ErrBadRequest")
)

// 500 series
var (
	ErrInternal           = XErr{Code: 500, Msg: "ErrInternal"}
	ErrRPCTimeout         = XErr{Code: 504, Msg: "rpc timeout"}
	ErrServiceUnavailable = XErr{Code: 510, Msg: "service unavailable"}
)

// Customized errors
var (
	_                      = 0
	ErrBizTimeout          = XErr{Code: 1000, Msg: "ErrBizTimeout"}
	ErrThirdParty          = XErr{Code: 1001, Msg: "ErrThirdParty"}
	ErrInvalidRegisterInfo = XErr{Code: 1002, Msg: "ErrInvalidRegisterInfo"}
)
