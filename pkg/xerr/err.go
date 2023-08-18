package xerr

import (
	"encoding/json"
	"fmt"
)

type XErr interface {
	error
	NewMsg(msg string) XErr
	AppendMsg(msg string) XErr
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
func (t E) NewMsg(msg string) XErr {
	t.Msg = msg
	return t
}
func (t E) AppendMsg(msg string) XErr {
	t.Msg += " |---append>> " + msg
	return t
}
