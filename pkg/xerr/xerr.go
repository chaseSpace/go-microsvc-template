package xerr

import "fmt"

type XErr interface {
	error
	_signXErr()
}

type T struct {
	Code int32
	Msg  string
}

func (t T) _signXErr() {
	panic("no need to implement me")
}

var _ XErr = new(T)

func (t T) Error() string {
	return fmt.Sprintf("XErr: code=%d msg=%s", t.Code, t.Msg)
}
func (t T) WithMsg(msg string) T {
	t.Msg = msg
	return t
}

var (
	ErrParams   = T{Code: 400, Msg: "request params err"}
	ErrInternal = T{Code: 501, Msg: "server internal err"}
)
