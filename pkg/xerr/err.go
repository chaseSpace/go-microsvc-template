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

func New(msg string, code ...int32) XErr {
	cd := ErrInternal.Code
	if len(code) > 0 {
		cd = code[0]
	}
	return XErr{Code: cd, Msg: msg}
}

// FromErr from error type to XErr, that might be fail then nil returned
func FromErr(err error) (t XErr, ok bool) {
	if err == nil {
		return ErrNil, true
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
	return ErrInternal.New(err.Error())
}

func (t XErr) FlatMsg() string {
	return fmt.Sprintf("code:%d - msg:%s", t.Code, t.Msg)
}

func (t XErr) Error() string {
	return util.ToJsonStr(&t)
}

func (t XErr) New(msg string, args ...any) XErr {
	t.Msg = fmt.Sprintf(msg, args...)
	return t
}
func (t XErr) AppendMsg(msg string, args ...any) XErr {
	t.Msg += " ➜ " + fmt.Sprintf(msg, args...)
	return t
}

func (t XErr) Equal(err error) bool {
	if err == nil {
		return t.IsNil()
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

func (t XErr) IsInternal() bool {
	return t.Code >= 500 && t.Code < 600
}

// IsNil helper function to XErr.IsNil()
func IsNil(err error) bool {
	return ToXErr(err).IsNil()
}
