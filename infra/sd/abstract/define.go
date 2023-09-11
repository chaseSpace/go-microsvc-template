package abstract

import (
	"context"
	"fmt"
)

type ServiceDiscovery interface {
	Name() string
	Register(serviceName string, host string, port int, metadata map[string]string) error
	Deregister(service string) error
	Discover(ctx context.Context, serviceName string, block bool) ([]ServiceInstance, error)
}

// ServiceInstance 表示注册的单个实例
type ServiceInstance struct {
	ID       string
	Name     string
	IsUDP    bool
	Host     string
	Port     int
	Metadata map[string]string
}

func (s ServiceInstance) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type CtxDurKey struct{}
