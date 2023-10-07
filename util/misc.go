package util

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"net"
	"strings"
	"sync"
)

type TcpListenerFetcher struct {
	portMin, portMax int
	mem              map[int]int
}

func NewTcpListenerFetcher(portMin, portMax int) *TcpListenerFetcher {
	return &TcpListenerFetcher{portMin: portMin, portMax: portMax, mem: make(map[int]int)}
}

func (t *TcpListenerFetcher) Get() (lis net.Listener, port int, err error) {
	if t.portMin >= t.portMax {
		return nil, 0, errors.New("portMin must less than portMax")
	}
	loops := t.portMax - t.portMin + 1
	for i := 0; i < loops; i++ {
		port = RandIntRange(t.portMin, t.portMax, t.mem)
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			if strings.Contains(err.Error(), "already") {
				//println("continue", port)
				continue
			}
			return nil, 0, err
		}
		//println(111, port)
		return
	}
	return nil, 0, fmt.Errorf("failed, tried %d times", loops)
}

var (
	// unique-id, copy-friendly, sortable by gen time
	// see https://github.com/segmentio/ksuid
	__ksuid      = ksuid.New()
	__ksuidMutex = sync.Mutex{}
)

func NewKsuid() string {
	__ksuidMutex.Lock()
	__ksuid = __ksuid.Next()
	__ksuidMutex.Unlock()
	return __ksuid.String()
}

func GetOptArg[T any](a []T, def T) T {
	if len(a) == 0 {
		return def
	}
	return a[0]
}

type FuzzyCharTyp int8

const (
	FuzzyCharTypNone FuzzyCharTyp = iota
	FuzzyCharTypPhone
	FuzzyCharTypCitizenId // 身份证
)

// FuzzyChars 对字符串进行模糊处理
// example：
//
//		-- 123456 => 12**56
//		-- 15983882334 => 159****2334
//	    -- 440308198612183456 => 440308********3456
func FuzzyChars(src string, typ ...FuzzyCharTyp) string {
	tp := GetOptArg[FuzzyCharTyp](typ, FuzzyCharTypNone)
	_lenDiv3 := len(src) / 3
	if _lenDiv3 == 0 {
		return ""
	}
	tmp := []rune(src)
	start := 0
	end := 0
	switch tp {
	case FuzzyCharTypPhone:
		ss := strings.Split(src, "-")
		if len(ss) == 2 {
			if len(ss[1])/3 > 0 {
				_lenDiv3 = len(ss[1]) / 3
				start = len(ss[0]) + 1 + _lenDiv3
				if _lenDiv3*3 != len(src) {
					end = len(src) - _lenDiv3
				} else {
					end = len(ss[0]) + 1 + _lenDiv3*2
				}
			} else {
				tp = FuzzyCharTypNone
			}

			goto OUTOF_SWITCH
		}

		if len(src) == 11 {
			start = 3
			end = 7
		}
	case FuzzyCharTypCitizenId:
		if len(src) == 18 {
			start = 6
			end = 14
		}
	default:
		if tp != FuzzyCharTypNone {
			panic(fmt.Sprintf("unknown FuzzyCharTyp:%v", tp))
		}
	}

OUTOF_SWITCH:
	if tp == FuzzyCharTypNone {
		start = _lenDiv3
		if _lenDiv3*3 != len(src) {
			end = len(src) - _lenDiv3
		} else {
			end = _lenDiv3 * 2
		}
	}

	for i := start; i < end; i++ {
		tmp[i] = '*'
	}
	return string(tmp)
}
