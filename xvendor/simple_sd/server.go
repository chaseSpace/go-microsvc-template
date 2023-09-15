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

var Sd ServiceDiscovery

func Init() {
	Sd = &SimpleSd{registry: map[string]*Service{}}
}

func (s *SimpleSd) Name() string {
	return "SimpleSd"
}

var (
	ErrInstanceNotRegistered = errors.New("instance not register")
)

func (s *SimpleSd) Register(instance ServiceInstance) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	serv := s.registry[instance.Name]
	if serv == nil {
		serv = newService(instance.Name)
		s.registry[instance.Name] = serv
	}
	instance.registerAt = time.Now()
	err := serv.Add(&instance)
	if err == nil {
		serv.eventUpdate.Broadcast()
	}
	return err
}

func (s *SimpleSd) Deregister(service, id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if serv := s.registry[service]; serv != nil {
		serv.eventUpdate.L.Lock()
		defer func() {
			serv.eventUpdate.L.Unlock()
			serv.eventUpdate.Broadcast()
		}()
		return serv.Remove(id)
	}
	return errors.Wrap(ErrInstanceNotRegistered, fmt.Sprintf("service: %s, id: %s", service, id))
}

func (s *SimpleSd) Discovery(ctx context.Context, service string, lastHash string) (instances []ServiceInstance, currHash string, err error) {
	serv := s.getService(service)

	instances, currHash = s.getInstances(service)
	// when lastHash is empty, here is no wait-time to discovery then returns
	if lastHash == "" {
		return instances, currHash, nil
	}

	go func() {
		<-ctx.Done()
		serv.eventUpdate.Broadcast()
	}()

	serv.eventUpdate.L.Lock()
	defer serv.eventUpdate.L.Unlock()
	for {
		serv.eventUpdate.Wait()
		instances, currHash = s.getInstances(service)
		if currHash != lastHash || lastHash == "" {
			return
		}
		if ctx.Err() != nil { // reach to deadline
			return
		}
	}
}

func (s *SimpleSd) HealthCheck(service, id string) bool {
	return s.getService(service).containsInstance(id)
}

func (s *SimpleSd) getService(service string) *Service {
	s.mu.Lock()
	defer s.mu.Unlock()
	if serv := s.registry[service]; serv != nil {
		return serv
	} else {
		serv = newService(service)
		s.registry[service] = serv
		return serv
	}
}

func (s *SimpleSd) getInstances(service string) ([]ServiceInstance, string) {
	s.mu.RLock()
	serv := s.registry[service]
	s.mu.RUnlock()

	// service must be existed (created by caller)
	list, currHash := serv.InstanceList()
	return lo.Map(list, func(item *ServiceInstance, index int) ServiceInstance {
		return *item
	}), currHash
}

type Service struct {
	service string
	mu      sync.RWMutex
	smap    map[string]*ServiceInstance

	currHash    string
	quit        chan struct{}
	eventUpdate *sync.Cond
}

const HealthCheckInterval = time.Second * 5
const HealthCheckMaxFails = 1

func newService(service string) *Service {
	s := &Service{service: service, smap: make(map[string]*ServiceInstance), eventUpdate: sync.NewCond(&sync.Mutex{})}
	s.currHash = EmptyInstanceHash
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

func (s *Service) containsInstance(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.smap[id] != nil
}

func (s *Service) resetHash() {
	var __inst []ServiceInstance
	for _, instance := range s.smap {
		__inst = append(__inst, *instance)
	}
	s.currHash = CalInstanceHash(__inst)
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
	hash = s.currHash
	return
}
