package xerr

import (
	"encoding/json"
	"fmt"
	"microsvc/util"
	"strings"
)

type XErr struct {
	Code int32
	Msg  string
}

// FromErr from error type to XErr, that might be fail then nil returned
func FromErr(err error) (t XErr, ok bool) {
	if err == nil {
		return t, true
	}
	if t, ok = err.(XErr); ok {
		return
	}
	// cross service transfer
	return FromErrStr(err.Error())
}

func FromErrStr(s string) (t XErr, ok bool) {
	_ = json.Unmarshal([]byte(s), &t)
	return t, t.Code > 0 && t.Msg != ""
}

// ToXErr Convert error type to XErr, a non-nil value returned
func ToXErr(err error) XErr {
	if t, ok := FromErr(err); ok {
		return t
	}
	return ErrInternal.NewMsg(err.Error())
}

func (t XErr) FlatMsg() string {
	return fmt.Sprintf("code:%d - msg:%s", t.Code, t.Msg)
}

func (t XErr) Error() string {
	return util.ToJsonStr(&t)
}

func (t XErr) NewMsg(msg string, args ...any) XErr {
	t.Msg = fmt.Sprintf(msg, args...)
	return t
}
func (t XErr) AppendMsg(msg string, args ...any) XErr {
	t.Msg += " ➜ " + fmt.Sprintf(msg, args...)
	return t
}

func (t XErr) Equal(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromErr(err); ok {
		return e.Code == t.Code
	}
	return false
}

func (t XErr) DeepEqual(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromErr(err); ok {
		return e.Code == t.Code && e.Msg == t.Msg
	}
	return false
}

func (t XErr) Is(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromErr(err); ok {
		return e.Code == t.Code && strings.HasPrefix(e.Msg, t.Msg)
	}
	return false
}

func (t XErr) IsNil() bool {
	return t.Code == ErrNil.Code
}

// IsNil helper function to XErr.IsNil()
func IsNil(err error) bool {
	return ToXErr(err).IsNil()
}
