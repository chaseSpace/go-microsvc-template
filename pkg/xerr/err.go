package xerr

import (
	"encoding/json"
	"fmt"
	"microsvc/util"
	"strings"
)

//type XErr interface {
//	error
//	Code() int32
//	Msg() string
//	FlatMsg() string
//	NewMsg(msg string, args ...any) E
//	AppendMsg(msg string, args ...any) E
//	DeepEqual(err error) bool
//	Equal(err error) bool
//	Is(err error) bool
//}

type E struct {
	ECode int32
	EMsg  string
}

// FromErr from error type to XErr, that might be fail then nil returned
func FromErr(err error) (t E, ok bool) {
	if t, ok = err.(E); ok {
		return
	}
	// cross service transform
	return FromErrStr(err.Error())
}

func FromErrStr(s string) (t E, ok bool) {
	_ = json.Unmarshal([]byte(s), &t)
	return t, t.ECode > 0 && t.EMsg != ""
}

// ToXErr Convert error type to XErr, a non-nil value returned
func ToXErr(err error) E {
	if t, ok := FromErr(err); ok {
		return t
	}
	return ErrInternal.NewMsg(err.Error())
}

func (t E) FlatMsg() string {
	return fmt.Sprintf("code:%d - msg:%s", t.ECode, t.EMsg)
}

func (t E) Error() string {
	return util.ToJsonStr(&t)
}

func (t E) NewMsg(msg string, args ...any) E {
	t.EMsg = fmt.Sprintf(msg, args...)
	return t
}
func (t E) AppendMsg(msg string, args ...any) E {
	t.EMsg += " ➜ " + fmt.Sprintf(msg, args...)
	return t
}

func (t E) Equal(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromErr(err); ok {
		return e.ECode == t.ECode
	}
	return false
}

func (t E) DeepEqual(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromErr(err); ok {
		return e.ECode == t.ECode && e.EMsg == t.EMsg
	}
	return false
}

func (t E) Is(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := FromErr(err); ok {
		return e.ECode == t.ECode && strings.HasPrefix(e.EMsg, t.EMsg)
	}
	return false
}
