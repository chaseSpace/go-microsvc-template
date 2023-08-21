package consul

import (
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"microsvc/infra/svcdiscovery/define"
)

type consulSD struct {
	client *capi.Client
}

var _ define.ServiceDiscovery = new(consulSD)

func NewConsulSD() (*consulSD, error) {
	// 默认连接 Consul HTTP API Addr> 127.0.0.1:8500
	cfg := capi.DefaultConfig()
	//cfg.Address 可修改
	client, err := capi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consulSD{client: client}, nil
}

func (c consulSD) Register(serviceName string, address string, port int, metadata map[string]string) error {
	tcpAddr := fmt.Sprintf("%s:%d", address, port)
	params := &capi.AgentServiceRegistration{
		Name:    serviceName,
		Tags:    []string{"microsvc"},
		Port:    port,
		Address: address,
		Meta:    metadata,
		Check:   newHealthCheck("microsvc-"+serviceName+"-health", tcpAddr),
	}
	return c.client.Agent().ServiceRegister(params)
}

func (c consulSD) Deregister(serviceName string) error {
	return c.client.Agent().ServiceDeregister(serviceName)
}

func (c consulSD) Discover(serviceName string) ([]define.ServiceInstance, error) {
	c.client.Agent().Services()
}
