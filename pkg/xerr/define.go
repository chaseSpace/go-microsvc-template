package xerr

var (
	ErrParams   = E{ECode: 400, EMsg: "ErrParams"}
	ErrInternal = E{ECode: 501, EMsg: "ErrInternal"}
	ErrTimeout  = E{ECode: 1000, EMsg: "ErrTimeout"}
)
