package simple_sd

import (
	"crypto/md5"
	"fmt"
	"github.com/samber/lo"
	"net"
	"sort"
	"strings"
	"time"
)

func netDialTest(isUDP bool, addr string, retry int, intvl time.Duration) bool {
	network := "tcp"
	if isUDP {
		network = "udp"
	}
	for i := 0; i < retry; i++ {
		conn, err := net.DialTimeout(network, addr, time.Second*3)
		if err == nil {
			_ = conn.Close()
			return true
		}
		Sdlogger.Debug("dial tcp %s err:%v", addr, err)
		time.Sleep(intvl)
	}
	return false
}

var EmptyInstanceHash = ""

func init() {
	hasher := md5.New()
	hasher.Write([]byte("EmptyInstanceHash"))
	EmptyInstanceHash = fmt.Sprintf("%x", hasher.Sum(nil))
}

func CalInstanceHash(instances []ServiceInstance) string {
	if len(instances) == 0 {
		return EmptyInstanceHash
	}
	addrs := lo.Map(instances, func(item ServiceInstance, index int) string {
		return item.Addr()
	})
	sort.Strings(addrs)

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(addrs, "")))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
