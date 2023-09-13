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
	registry  map[string]string // svc -> id
}

var _ abstract.ServiceDiscovery = (*Consul)(nil)

const Name = "Consul"

func New() (*Consul, error) {
	cfg := capi.DefaultConfig()
	client, err := capi.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Consul{client: client}, nil
}

func (c *Consul) Name() string {
	return Name
}

func (c *Consul) Register(serviceName string, host string, port int, metadata map[string]string) error {
	if c.registry[serviceName] != "" {
		return fmt.Errorf("consul: already registered")
	}
	id := util.RandomString(4)
	tcpAddr := fmt.Sprintf("%s:%d", host, port)
	params := &capi.AgentServiceRegistration{
		ID:      id,
		Name:    serviceName,
		Tags:    []string{"microsvc"},
		Port:    port,
		Address: host,
		Meta:    metadata,
		Check:   newHealthCheck("microsvc-"+serviceName+"-health", tcpAddr),
	}

	err := c.client.Agent().ServiceRegister(params)
	if err != nil {
		return err
	}
	c.registry[serviceName] = id
	return nil
}

func (c *Consul) Deregister(service string) error {
	if c.registry[service] == "" {
		return fmt.Errorf("consul: not register")
	}
	delete(c.registry, service)
	return c.client.Agent().ServiceDeregister(service)
}

func (c *Consul) Discover(ctx context.Context, serviceName string, block bool) (list []abstract.ServiceInstance, err error) {
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
			Host:     s.Service.Address,
			Port:     s.Service.Port,
			Metadata: s.Service.Meta,
		}
		list = append(list, inst)
	}
	return list, nil
}
