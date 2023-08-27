package util

import (
	"net"
	"testing"
)

var ss []net.Listener

func closeAll() {
	for _, s := range ss {
		s.Close()
	}
}

func TestGetTcpListenerWithinRangePort(t *testing.T) {
	// success case
	start := 1000
	length := 100
	f := NewTcpListenerFetcher(start, start+length)
	for i := 0; i < length; i++ {
		lis, port, err := f.Get()
		AssertNil(err)
		AssertNotNil(lis)
		if !(start <= port && port <= start+length) {
			t.Fatalf("case 1 - err port:%v", port)
		}
		println("ok", port)
		ss = append(ss, lis)
	}
	closeAll()
}
