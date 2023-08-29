package xerr

// Common error
var (
	_           = 0
	ErrOK       = E{ECode: 200, EMsg: "OK"}
	ErrParams   = E{ECode: 400, EMsg: "ErrParams"}
	ErrInternal = E{ECode: 501, EMsg: "ErrInternal"}
	ErrTimeout  = E{ECode: 1000, EMsg: "ErrTimeout"}
)

// Internal error
var (
	_              = 0
	ErrNoRPCClient = ErrInternal.NewMsg("no available rpc client")
)
