package simple_sd

import (
	"context"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"microsvc/infra/sd/abstract"
	"microsvc/pkg/xerr"
	"microsvc/util"
	"time"
)

type SimpleSd struct {
	serverPort int
	lastHash   string
	registry   map[string]string // svc -> id
}

func New(port int) *SimpleSd {
	return &SimpleSd{serverPort: port, registry: make(map[string]string)}
}

var _ abstract.ServiceDiscovery = (*SimpleSd)(nil)

const (
	Name          = "simple_sd"
	httpResOkCode = 200

	registerPath   = "/service/register"
	deregisterPath = "/service/deregister"
	discoveryPath  = "/service/discovery"
)

func (s *SimpleSd) getRequestUrl(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", s.serverPort, path)
}

func (s *SimpleSd) Name() string {
	return Name
}

type httpRes struct {
	Code int // 200 OK
	Msg  string
	Data interface{} `json:"Data,omit_empty"`
}

func (s *SimpleSd) Register(service string, host string, port int, metadata map[string]string) error {
	if s.registry[service] != "" {
		return fmt.Errorf("already registered")
	}
	type registerReq struct {
		abstract.ServiceInstance
	}
	req := &registerReq{abstract.ServiceInstance{
		ID:       util.RandomString(4),
		Name:     service,
		IsUDP:    false,
		Host:     host,
		Port:     port,
		Metadata: metadata,
	}}
	res := new(httpRes)
	_, _, errs := gorequest.New().Post(s.getRequestUrl(registerPath)).SendStruct(req).EndStruct(res)
	if len(errs) > 0 {
		return errs[0]
	}
	if res.Code != httpResOkCode {
		return xerr.ErrInternal.NewMsg("register failed, got resp: %+v", res)
	}
	s.registry[service] = req.ID
	return nil
}

func (s *SimpleSd) Deregister(service string) error {
	id := s.registry[service]
	if id == "" {
		return xerr.ErrInternal.NewMsg("simpleSd: not register")
	}
	type deregisterReq struct {
		Service string
		Id      string
	}
	req := &deregisterReq{
		Service: service,
		Id:      id,
	}
	res := new(httpRes)
	_, _, errs := gorequest.New().Post(s.getRequestUrl(deregisterPath)).SendStruct(req).EndStruct(res)
	if len(errs) > 0 {
		return errs[0]
	}
	if res.Code != httpResOkCode {
		return xerr.ErrInternal.NewMsg("deregister failed, got resp: %+v", res)
	}
	delete(s.registry, id)
	return nil
}

func (s *SimpleSd) Discover(ctx context.Context, serviceName string, block bool) ([]abstract.ServiceInstance, error) {
	type discoveryReq struct {
		Service   string
		LastHash  string
		WaitMaxMs int64
	}
	type discoveryRsp struct {
		Instances []abstract.ServiceInstance
		Hash      string
	}
	req := &discoveryReq{
		Service:   serviceName,
		LastHash:  s.lastHash,
		WaitMaxMs: time.Minute.Milliseconds() * 2,
	}
	if !block {
		req.LastHash = ""
	}
	//println(111, req.LastHash)
	data := &discoveryRsp{}
	res := &httpRes{Data: data}
	_, _, errs := gorequest.New().Post(s.getRequestUrl(discoveryPath)).SendStruct(req).EndStruct(res)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if res.Code != httpResOkCode {
		return nil, xerr.ErrInternal.NewMsg("discovery failed, got resp: %+v", res)
	}
	s.lastHash = data.Hash
	return data.Instances, nil
}
