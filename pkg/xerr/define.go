package xerr

var (
	ErrParams   = E{Ecode: 400, Msg: "request params err"}
	ErrInternal = E{Ecode: 501, Msg: "server internal err"}
	ErrTimeout  = E{Ecode: 1000, Msg: "operation timeout"}
)
