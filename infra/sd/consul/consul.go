package consul

import (
	"context"
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"microsvc/infra/sd/abstract"
	"microsvc/util"
	"time"
)

type ConsulSD struct {
	client    *capi.Client
	lastIndex uint64
}

var _ abstract.ServiceDiscovery = (*ConsulSD)(nil)

func New() (*ConsulSD, error) {
	// 默认连接 Consul HTTP API Addr> 127.0.0.1:8500
	cfg := capi.DefaultConfig()
	//cfg.Address 可修改
	client, err := capi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ConsulSD{client: client}, nil
}

func (c *ConsulSD) Register(serviceName string, address string, port int, metadata map[string]string) error {
	tcpAddr := fmt.Sprintf("%s:%d", address, port)
	params := &capi.AgentServiceRegistration{
		//addr:  默认等于Name
		Name:    serviceName,
		Tags:    []string{"microsvc"},
		Port:    port,
		Address: address,
		Meta:    metadata,
		Check:   newHealthCheck("microsvc-"+serviceName+"-health", tcpAddr),
	}

	err := c.client.Agent().ServiceRegister(params)
	return err
}

func (c *ConsulSD) Deregister(serviceName string) error {
	return c.client.Agent().ServiceDeregister(serviceName)
}

func (c *ConsulSD) Discover(ctx context.Context, serviceName string) (list []abstract.ServiceInstance, err error) {
	err = context.DeadlineExceeded // default
	dur := time.Minute
	if val := ctx.Value(abstract.CtxDurKey{}); val != nil {
		dur = val.(time.Duration)
	}
	util.RunTask(ctx, func() {
		list, err = c.getInstances(serviceName, dur)
	})
	return
}

func (c *ConsulSD) getInstances(serviceName string, waitTime time.Duration) (list []abstract.ServiceInstance, err error) {
	opt := &capi.QueryOptions{WaitIndex: c.lastIndex, WaitTime: waitTime}
	entries, meta, err := c.client.Health().Service(serviceName, "", true, opt)
	if err != nil {
		return nil, err
	}
	if c.lastIndex > meta.LastIndex { //  index goes backwards, reset it
		c.lastIndex = 0
	} else if c.lastIndex < meta.LastIndex {
		c.lastIndex = meta.LastIndex
	}
	for _, s := range entries {
		inst := abstract.ServiceInstance{
			ID:       s.Service.ID,
			Name:     serviceName,
			Address:  s.Service.Address,
			Port:     s.Service.Port,
			Metadata: s.Service.Meta,
		}
		list = append(list, inst)
	}
	return list, nil
}
