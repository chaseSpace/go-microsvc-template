package mdns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"microsvc/infra/sd/abstract"
	"microsvc/util"
	"microsvc/xvendor/mdns"
	"net"
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
	registry  map[string]*registry // svc -> id
	instCache map[string]map[string]int8
}

type registry struct {
	s  *mdns.Server
	id string
}

var _ abstract.ServiceDiscovery = (*Mdns)(nil)

// mDnsDomain is the domain name used by mDNS. this is necessary to
// set, because mDNS is a multicast protocol and the domain name
// is used to filter the packets.
const mDnsDomain = "microsvc."
const Name = "mDNS"

func New() *Mdns {
	return &Mdns{registry: make(map[string]*registry), instCache: map[string]map[string]int8{}}
}

func (m *Mdns) Name() string {
	return Name
}

func getServerId(svc, host string, port int) string {
	return fmt.Sprintf("%s:%s:%d", svc, host, port)
}

func (m *Mdns) Register(serviceName string, host string, port int, metadata map[string]string) (err error) {
	if m.registry[serviceName] != nil {
		return fmt.Errorf("mdns: already registered")
	}
	id := util.RandomString(4)
	s, err := mdns.NewMDNSService(id, serviceName, mDnsDomain, "", port,
		[]net.IP{[]byte{127, 0, 0, 1}}, []string{util.ToJsonStr(metadata)})
	if err != nil {
		return err
	}
	server, err := mdns.NewServer(&mdns.Config{Zone: s})
	if err != nil {
		return err
	}
	m.registry[serviceName] = &registry{
		s:  server,
		id: id,
	}
	return
}

func (m *Mdns) Deregister(service string) error {
	if rs := m.registry[service]; rs == nil {
		return fmt.Errorf("mdns: not register")
	} else {
		delete(m.registry, service)
		return rs.s.Shutdown()
	}
}

func (m *Mdns) Discover(ctx context.Context, svc string, block bool) (instances []abstract.ServiceInstance, err error) {
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
				md := make(map[string]string)
				_ = json.Unmarshal([]byte(entry.Info), &md)
				instances = append(instances, abstract.ServiceInstance{
					Name:     entry.Name, // hostname.service.domain by mDNS
					Host:     entry.AddrV4.String(),
					Port:     entry.Port,
					Metadata: md,
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
