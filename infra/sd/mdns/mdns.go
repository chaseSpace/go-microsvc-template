package mdns

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"microsvc/infra/sd/abstract"
	"microsvc/util"
	"microsvc/xvendor/mdns"
	"net"
	"os"
	"time"
)

/*
INTRODUCTION:
 	mDNS is a multicast protocol that was developed by Apple Inc, and
	the domain name is used to filter the packets.

WARNING:
	mDNS is not suitable for production environments. It is designed
	for local networks only.
	And, In Windows, mDNS support may not be stable by default. To
	completely support mDNS functionality, it is recommended to install
	Bonjour Print Services for Windows.
*/

// Mdns implements the abstract.ServiceDiscovery with mDNS (Multicast DNS)
// protocol using UDP.
// Note: Mdns should not be used in a production environment.
type Mdns struct {
	server    map[string]*mdns.Server
	instCache map[string]map[string]int8
}

var _ abstract.ServiceDiscovery = (*Mdns)(nil)

// mDnsDomain is the domain name used by mDNS. this is necessary to
// set, because mDNS is a multicast protocol and the domain name
// is used to filter the packets.
const mDnsDomain = "microsvc."
const Name = "mDNS"

func New() *Mdns {
	return &Mdns{server: make(map[string]*mdns.Server), instCache: map[string]map[string]int8{}}
}

func (m *Mdns) Name() string {
	return Name
}

func (m *Mdns) Register(serviceName string, address string, port int, metadata map[string]string) (err error) {
	if m.server[serviceName] != nil {
		return fmt.Errorf("already registered")
	}
	host, _ := os.Hostname()
	s, err := mdns.NewMDNSService(host+util.RandomString(3), serviceName, mDnsDomain, "", port,
		[]net.IP{[]byte{127, 0, 0, 1}}, nil)
	if err != nil {
		return err
	}
	server, err := mdns.NewServer(&mdns.Config{Zone: s})
	if err != nil {
		return err
	}
	m.server[serviceName] = server
	return
}

func (m *Mdns) Deregister(serviceName string) error {
	if ser := m.server[serviceName]; ser != nil {
		return ser.Shutdown()
	}
	return nil
}

func (m *Mdns) Discovery(ctx context.Context, svc string, block bool) (instances []abstract.ServiceInstance, err error) {
	asyncRecv := func(entries chan *mdns.ServiceEntry) {
		defer func() {
			//fmt.Printf("2222 %+v  %s\n", instances, svc)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case entry := <-entries:
				if entry == nil { // channel closed
					return
				}
				instances = append(instances, abstract.ServiceInstance{
					Name:     entry.Name, // hostname.service.domain by mDNS
					Address:  entry.AddrV4.String(),
					Port:     entry.Port,
					Metadata: nil,
				})
			}
		}
	}

	for {
		entries := make(chan *mdns.ServiceEntry, 4)
		go asyncRecv(entries)

		err = mdns.Lookup(svc, mDnsDomain, entries)
		close(entries)
		if err != nil {
			return
		}

		_, changed := m.updateCache(svc, instances)
		if changed || !block {
			return
		}

		instances = nil
		time.Sleep(time.Second * 5) // could tune-up, 5~30s is a good choice
	}
}

func (m *Mdns) updateCache(serviceName string, instances []abstract.ServiceInstance) (map[string]int8, bool) {
	changed := false

	var newCache map[string]int8
	cache := m.instCache[serviceName]
	if cache == nil {
		cache = make(map[string]int8)
		newCache = cache
		m.instCache[serviceName] = cache
		changed = true
	} else {
		newCache = make(map[string]int8)
	}

	lo.ForEach(instances, func(item abstract.ServiceInstance, index int) {
		if cache[item.Addr()] == 0 {
			changed = true
		}
		newCache[item.Addr()] = 1
	})

	m.instCache[serviceName] = newCache

	return cache, changed || len(instances) != len(cache)
}
