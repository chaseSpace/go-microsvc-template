package abstract

import (
	"context"
	"fmt"
)

type ServiceDiscovery interface {
	Register(serviceName string, address string, port int, metadata map[string]string) error
	Deregister(serviceName string) error
	// Discover must support block query
	Discovery(ctx context.Context, serviceName string) ([]ServiceInstance, error)
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
