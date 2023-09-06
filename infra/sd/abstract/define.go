package abstract

import (
	"context"
	"fmt"
)

type ServiceDiscovery interface {
	Name() string
	Register(serviceName string, address string, port int, metadata map[string]string) error
	Deregister(serviceName string) error
	Discovery(ctx context.Context, serviceName string, block bool) ([]ServiceInstance, error)
}

// ServiceInstance 表示注册的单个实例
type ServiceInstance struct {
	ID       string
	Name     string
	Address  string
	Port     int
	Metadata map[string]string
}

func (s ServiceInstance) Addr() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

type CtxDurKey struct{}
