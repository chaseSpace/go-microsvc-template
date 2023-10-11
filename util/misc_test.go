package util

import (
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
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

func TestFuzzyChars(t *testing.T) {
	// success case
	assert.Equal(t, "h*l", FuzzyChars("hel"))
	assert.Equal(t, "h***o", FuzzyChars("hello"))
	assert.Equal(t, "he**o1", FuzzyChars("hello1"))

	assert.Equal(t, "158****8888", FuzzyChars("15899998888", FuzzyCharTypPhone))
	assert.Equal(t, "020-29***23", FuzzyChars("020-2938123", FuzzyCharTypPhone))
	assert.Equal(t, "0***2", FuzzyChars("020-2", FuzzyCharTypPhone))
	assert.Equal(t, "020-2*4", FuzzyChars("020-234", FuzzyCharTypPhone))

	assert.Equal(t, "440308********3456", FuzzyChars("440308198612183456", FuzzyCharTypCitizenId))
}

func TestNewKsuid(t *testing.T) {
	var ids sync.Map

	var x sync.WaitGroup

	loops := 1000
	for i := 0; i < loops; i++ {
		x.Add(1)
		go func() {
			ids.Store(NewKsuid(), 1)
			x.Done()
		}()
	}
	x.Wait()

	actuals := 0
	ids.Range(func(key, value any) bool {
		actuals++
		return true
	})
	assert.Equal(t, loops, actuals)
}
