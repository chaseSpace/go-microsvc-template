package util

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
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
