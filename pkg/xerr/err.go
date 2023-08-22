package xerr

import (
	"encoding/json"
	"fmt"
)

type XErr interface {
	error
	NewMsg(msg string, args ...any) XErr
	AppendMsg(msg string, args ...any) XErr
}

type E struct {
	Ecode int32
	Msg   string
}

func IsXErr(err error) (bool, XErr) {
	t := new(E)
	_ = json.Unmarshal([]byte(err.Error()), t)
	if t.Ecode > 0 {
		return true, t
	}
	return false, nil
}

var _ XErr = new(E)

func (t E) Error() string {
	return fmt.Sprintf("XERR: ecode=%d msg=%s", t.Ecode, t.Msg)
}
func (t E) NewMsg(msg string, args ...any) XErr {
	t.Msg = fmt.Sprintf(msg, args...)
	return t
}
func (t E) AppendMsg(msg string, args ...any) XErr {
	t.Msg += " |---append>> " + fmt.Sprintf(msg, args...)
	return t
}
