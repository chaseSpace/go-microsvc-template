package simple_sd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"sort"
	"sync"
	"time"
)

type SimpleSd struct {
	mu               sync.RWMutex
	registry         map[string]*Service
	newServiceNotify sync.Cond
}

var Sd ServiceDiscovery = &SimpleSd{registry: map[string]*Service{}}

func (s *SimpleSd) Name() string {
	return "SimpleSd"
}

var (
	ErrInstanceNotRegistered = errors.New("instance not register")
)

const DiscoveryInterval = time.Millisecond * 200

func (s *SimpleSd) Register(instance ServiceInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	serv := s.registry[instance.Service]
	if serv == nil {
		serv = newService(instance.Service)
		s.registry[instance.Service] = serv
	}
	instance.registerAt = time.Now()
	err := serv.Add(&instance)
	return err
}

func (s *SimpleSd) Deregister(service, id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if serv := s.registry[service]; serv != nil {
		return serv.Remove(id)
	}
	return errors.Wrap(ErrInstanceNotRegistered, fmt.Sprintf("service: %s, id: %s", service, id))
}

func (s *SimpleSd) Discovery(ctx context.Context, service string, lastHash string) ([]ServiceInstance, string, error) {
	// when lastHash is empty, here is realtime discovery and returns
	if lastHash == "" {
		instances, currHash := s.getInstances(service)
		return instances, currHash, nil
	}
	for {
		select {
		case <-ctx.Done():
			instances, currHash := s.getInstances(service)
			return instances, currHash, nil

		case <-time.After(DiscoveryInterval): // refresh interval
			instances, currHash := s.getInstances(service)
			if currHash != lastHash || lastHash == "" {
				return instances, currHash, nil
			}
		}
	}
}

func (s *SimpleSd) getInstances(service string) ([]ServiceInstance, string) {
	s.mu.RLock()
	serv := s.registry[service]
	s.mu.RUnlock()

	if serv == nil {
		return nil, ""
	} else {
		list, currHash := serv.InstanceList()
		return lo.Map(list, func(item *ServiceInstance, index int) ServiceInstance {
			return *item
		}), currHash
	}
}

type Service struct {
	service string
	mu      sync.RWMutex
	smap    map[string]*ServiceInstance

	nowHash string
	quit    chan struct{}
}

const HealthCheckInterval = time.Second * 5
const HealthCheckMaxFails = 1

func newService(service string) *Service {
	s := &Service{service: service, smap: make(map[string]*ServiceInstance)}
	go s.healthCheck()
	return s
}

func (s *Service) Stop() {
	close(s.quit)
}

func (s *Service) healthCheck() {
	Sdlogger.Debug("service %s launched health check", s.service)
	for {
		select {
		case <-s.quit:
			Sdlogger.Debug("service: %s health check stopped", s.service)
			return
		case <-time.After(HealthCheckInterval):
			__inst, _ := s.InstanceList()
			for _, ins := range __inst {
				st := time.Now()
				pass := netDialTest(ins.IsUDP, ins.Addr(), 1, time.Millisecond*20)
				cost := time.Since(st)
				if pass {
					if ins.fails > 0 {
						ins.fails = 0
						Sdlogger.Debug("service: %s, instance %s health recovered", s.service, ins.Addr())
					}
				} else {
					println(2222, ins.fails, cost.String())
					ins.fails++
					if ins.fails >= HealthCheckMaxFails {
						// remove instance
						_ = s.Remove(ins.Addr())
						Sdlogger.Debug("service: %s, instance %s was removed", s.service, ins.Addr())
					}
				}
			}
		}
	}
}

func (s *Service) Add(instance *ServiceInstance) error {
	s.mu.Lock()
	s.smap[instance.Id] = instance
	s.mu.Unlock()
	s.resetHash()
	return nil
}

func (s *Service) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.smap[id] == nil {
		return fmt.Errorf("instance: %s not found", id)
	}
	delete(s.smap, id)
	s.resetHash()
	return nil
}

func (s *Service) resetHash() {
	var __inst []ServiceInstance
	for _, instance := range s.smap {
		__inst = append(__inst, *instance)
	}
	s.nowHash = CalInstanceHash(__inst)
}

func (s *Service) InstanceList() (list []*ServiceInstance, hash string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, instance := range s.smap {
		list = append(list, instance)
	}
	// sort in ascending order by registration time
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].registerAt.Before(list[j].registerAt)
	})
	hash = s.nowHash
	return
}
