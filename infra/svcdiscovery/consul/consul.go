package consul

import (
	"context"
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"microsvc/infra/svcdiscovery/define"
)

type ConsulSD struct {
	client *capi.Client
}

var _ define.ServiceDiscovery = new(ConsulSD)

func NewConsulSD() (*ConsulSD, error) {
	// 默认连接 Consul HTTP API Addr> 127.0.0.1:8500
	cfg := capi.DefaultConfig()
	//cfg.Address 可修改
	client, err := capi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ConsulSD{client: client}, nil
}

func (c ConsulSD) Register(serviceName string, address string, port int, metadata map[string]string) error {
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

func (c ConsulSD) Deregister(serviceName string) error {
	return c.client.Agent().ServiceDeregister(serviceName)
}

func (c ConsulSD) Discover(ctx context.Context, serviceName string) (inst []define.ServiceInstance, err error) {
	return
}

func (c ConsulSD) getInstances(serviceName string, waitHash string) (list []define.ServiceInstance, lastHash string, err error) {

	opt := &capi.QueryOptions{WaitHash: waitHash}
	entries, meta, err := c.client.Health().Service(serviceName, "", true, opt)
	if err != nil {
		return nil, "", err
	}

	for _, s := range entries {
		inst := define.ServiceInstance{
			ID:       s.Service.ID,
			Name:     serviceName,
			Address:  s.Service.Address,
			Port:     s.Service.Port,
			Metadata: s.Service.Meta,
		}
		list = append(list, inst)
	}
	return list, meta.LastContentHash, nil
}
