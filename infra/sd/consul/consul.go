package consul

import (
	"context"
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"microsvc/infra/sd/abstract"
	"microsvc/util"
	"time"
)

type Consul struct {
	client    *capi.Client
	lastIndex uint64
}

var _ abstract.ServiceDiscovery = (*Consul)(nil)

const Name = "Consul"

func New() (*Consul, error) {
	// 默认连接 Consul HTTP API Addr> 127.0.0.1:8500
	cfg := capi.DefaultConfig()
	//cfg.Address 可修改
	client, err := capi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Consul{client: client}, nil
}

func (c *Consul) Name() string {
	return Name
}

func (c *Consul) Register(serviceName string, address string, port int, metadata map[string]string) error {
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

func (c *Consul) Deregister(serviceName string) error {
	return c.client.Agent().ServiceDeregister(serviceName)
}

func (c *Consul) Discovery(ctx context.Context, serviceName string, block bool) (list []abstract.ServiceInstance, err error) {
	err = context.DeadlineExceeded // default
	dur := time.Minute
	if val := ctx.Value(abstract.CtxDurKey{}); val != nil {
		dur = val.(time.Duration) // use duration here, because Consul do not support block by context
	}
	util.RunTask(ctx, func() {
		list, err = c.getInstances(serviceName, dur, block)
	})
	return
}

func (c *Consul) getInstances(serviceName string, waitTime time.Duration, block bool) (list []abstract.ServiceInstance, err error) {
	opt := &capi.QueryOptions{WaitIndex: c.lastIndex, WaitTime: waitTime}
	if !block {
		opt.WaitIndex = 0 // set to 0 to disable blocking query
	}
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
