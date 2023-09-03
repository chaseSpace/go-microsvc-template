package mdns

import (
	"context"
	"fmt"
	"github.com/hashicorp/mdns"
	"microsvc/infra/sd/abstract"
	"os"
)

// Mdns implements the abstract.ServiceDiscovery with mDNS (Multicast DNS) protocol using UDP.
// Note: Mdns should not be used in a production environment.
type Mdns struct {
	server map[string]*mdns.Server
}

var _ abstract.ServiceDiscovery = (*Mdns)(nil)

func New() *Mdns {
	return &Mdns{server: make(map[string]*mdns.Server)}
}

func (m Mdns) Register(serviceName string, address string, port int, metadata map[string]string) (err error) {
	if m.server[serviceName] != nil {
		return fmt.Errorf("already registered")
	}
	host, _ := os.Hostname()
	println(1111, host)
	s, _ := mdns.NewMDNSService(serviceName+"."+host, serviceName, "", address, port, nil, nil)

	server, err := mdns.NewServer(&mdns.Config{Zone: s})
	if err != nil {
		return err
	}
	m.server[serviceName] = server
	return err
}

func (m Mdns) Deregister(serviceName string) error {
	if ser := m.server[serviceName]; ser != nil {
		return ser.Shutdown()
	}
	return nil
}

func (m Mdns) Discover(ctx context.Context, serviceName string) (insts []abstract.ServiceInstance, err error) {
	entries := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entries {
			insts = append(insts, abstract.ServiceInstance{
				Name:     entry.Name,
				Address:  entry.Host,
				Port:     entry.Port,
				Metadata: nil,
			})
		}
	}()
	_ = mdns.Lookup(serviceName, entries)
	close(entries)
	return
}
