package abstract

import (
	"context"
	"time"
)

type EmptySD struct {
}

var _ ServiceDiscovery = new(EmptySD)

func (e EmptySD) Register(serviceName string, address string, port int, metadata map[string]string) error {
	return nil
}

func (e EmptySD) Deregister(serviceName string) error {
	return nil
}

func (e EmptySD) Discover(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	time.Sleep(time.Millisecond * 100) // mock block
	return nil, nil
}
