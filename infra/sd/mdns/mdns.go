package mdns

import (
	"context"
	"fmt"
	"github.com/hashicorp/mdns"
	"github.com/samber/lo"
	"microsvc/infra/sd/abstract"
	"microsvc/util"
	"os"
	"time"
)

// MdnsSD implements the abstract.ServiceDiscovery with mDNS (Multicast DNS) protocol using UDP.
// Note: MdnsSD should not be used in a production environment.
type MdnsSD struct {
	server    map[string]*mdns.Server
	instCache map[string]map[string]int8
}

var _ abstract.ServiceDiscovery = (*MdnsSD)(nil)

const domain = "microsvc."

func New() *MdnsSD {
	return &MdnsSD{server: make(map[string]*mdns.Server), instCache: map[string]map[string]int8{}}
}

func (m *MdnsSD) Register(serviceName string, address string, port int, metadata map[string]string) (err error) {
	if m.server[serviceName] != nil {
		return fmt.Errorf("already registered")
	}
	host, _ := os.Hostname()
	s, err := mdns.NewMDNSService(host+util.RandomString(3), serviceName, domain, "", port, nil, nil)
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

func (m *MdnsSD) Deregister(serviceName string) error {
	if ser := m.server[serviceName]; ser != nil {
		return ser.Shutdown()
	}
	return nil
}

func (m *MdnsSD) Discovery(ctx context.Context, serviceName string) (instances []abstract.ServiceInstance, err error) {
	asyncRecv := func(entries chan *mdns.ServiceEntry) {
		for entry := range entries {
			instances = append(instances, abstract.ServiceInstance{
				Name:     entry.Name, // hostname.service.domain by mDNS
				Address:  "127.0.0.1",
				Port:     entry.Port,
				Metadata: nil,
			})
		}
		fmt.Printf("2222  %+v  \n", instances)
	}

	// block query until ctx timeout
	for {
		select {
		case <-ctx.Done():
			return
		default:
			entries := make(chan *mdns.ServiceEntry, 4)
			go asyncRecv(entries)

			p := mdns.DefaultParams(serviceName)
			p.DisableIPv6 = true
			p.Entries = entries
			p.Domain = domain
			err = mdns.Query(p)
			close(entries)
			if err != nil {
				return
			}

			_, updated := m.updateCache(serviceName, instances)
			if updated {
				return
			}

			instances = nil
			time.Sleep(time.Second * 2) // could tune-up
		}
	}
}

func (m *MdnsSD) updateCache(serviceName string, instances []abstract.ServiceInstance) (map[string]int8, bool) {
	updated := false

	var newCache map[string]int8
	cache := m.instCache[serviceName]
	if cache == nil {
		cache = make(map[string]int8)
		newCache = cache
		m.instCache[serviceName] = cache
		updated = true
	} else {
		newCache = make(map[string]int8)
	}

	lo.ForEach(instances, func(item abstract.ServiceInstance, index int) {
		if cache[item.Addr()] == 0 {
			updated = true
		}
		newCache[item.Addr()] = 1
	})

	m.instCache[serviceName] = newCache

	if !updated {
		updated = len(instances) != len(cache)
	}

	return cache, updated
}
