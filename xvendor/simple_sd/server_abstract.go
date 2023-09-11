package simple_sd

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ServiceDiscovery interface {
	Name() string
	Register(instance ServiceInstance) error
	Deregister(service, id string) error
	Discovery(ctx context.Context, service string, lastHash string) ([]ServiceInstance, string, error)
}

// ServiceInstance 表示注册的单个实例
type ServiceInstance struct {
	Id       string
	Service  string
	IsUDP    bool // TCP by default
	Host     string
	Port     int
	Metadata map[string]string

	registerAt time.Time
	fails      int
}

func (s ServiceInstance) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s ServiceInstance) Check() error {
	if s.Service == "" || s.Id == "" {
		return errors.New("ServiceInstance must have service name and id")
	}
	if s.Host == "" || s.Port < 1 {
		return errors.New("ServiceInstance must have valid address and port")
	}
	return nil
}
