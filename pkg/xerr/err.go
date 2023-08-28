package xerr

import (
	"encoding/json"
	"fmt"
	"microsvc/util"
)

type XErr interface {
	error
	Code() int32
	Msg() string
	NewMsg(msg string, args ...any) XErr
	AppendMsg(msg string, args ...any) XErr
}

type E struct {
	ECode int32
	EMsg  string
}

// FromErr from error type to XErr, that might be fail then nil returned
func FromErr(err error) XErr {
	if e, _ := err.(XErr); e != nil {
		return e
	}
	// cross service transform
	t := new(E)
	_ = json.Unmarshal([]byte(err.Error()), t)
	if t.ECode > 0 {
		return t
	}
	return nil
}

// ToXErr Convert error type to XErr, a non-nil value returned
func ToXErr(err error) XErr {
	if e := FromErr(err); e != nil {
		return e
	}
	return ErrInternal.NewMsg(err.Error())
}

var _ XErr = new(E)

func (t E) Code() int32 {
	return t.ECode
}

func (t E) Msg() string {
	return t.EMsg
}

func (t E) Error() string {
	return util.ToJsonStr(&t)
}

func (t E) NewMsg(msg string, args ...any) XErr {
	t.EMsg = fmt.Sprintf(msg, args...)
	return t
}
func (t E) AppendMsg(msg string, args ...any) XErr {
	t.EMsg += " |err.append>> " + fmt.Sprintf(msg, args...)
	return t
}

func (t E) Equal(err error) bool {
	if err == nil {
		return false
	}
	if e := FromErr(err); e == nil {
		return false
	} else {
		return e.Code() == t.ECode && e.Msg() == t.EMsg
	}
}
