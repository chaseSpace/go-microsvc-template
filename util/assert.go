package util

import "fmt"

func AssertNilErr(v error) {
	if v != nil {
		panic(fmt.Sprintf("AssertNilErr: non-nil err:%v", v))
	}
}

func AssertNil(v interface{}) {
	if v != nil {
		panic(fmt.Sprintf("AssertNil: non-nil val:%+v", v))
	}
}

func AssertNotNil(v interface{}) {
	if v == nil {
		panic(fmt.Sprintf("AssertNil: nil val, type:%T", v))
	}
}
