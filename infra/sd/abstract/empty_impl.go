package abstract

import (
	"context"
	"time"
)

type Empty struct {
}

var _ ServiceDiscovery = new(Empty)

func (e Empty) Name() string {
	return "empty"
}

func (e Empty) Register(serviceName string, address string, port int, metadata map[string]string) error {
	return nil
}

func (e Empty) Deregister(serviceName string) error {
	return nil
}

func (e Empty) Discover(ctx context.Context, serviceName string, block bool) ([]ServiceInstance, error) {
	if block {
		time.Sleep(time.Millisecond * 100) // mock block
	}
	return nil, nil
}

func (e Empty) HealthCheck(ctx context.Context, service string) error {
	return nil
}
