package simple_sd

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type HttpRes struct {
	Code int // 200 OK
	Msg  string
	Data interface{} `json:"Data,omit_empty"`
}

func ToJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

var (
	ErrMethod       = errors.New("use HTTP POST")
	ErrParams       = errors.New("invalid params")
	ErrInstanceAddr = errors.New("invalid Addr, provide a addr like localhost:8000 or 192.168.1.1:8000")
)

func newRes(data interface{}, code int, err error) *HttpRes {
	msg := "OK"
	if err != nil {
		msg = err.Error()
	}
	return &HttpRes{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

type registerReq struct {
	ServiceInstance
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(registerReq)

	var data []ServiceInstance
	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(data, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleRegister: service:%s req:%s error: %v", req.Service, req.Addr(), err)
			return
		}
		Sdlogger.Info("handleRegister OK, service:%s req:%s", req.Service, req.Addr())
	}()

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()

	if err = req.Check(); err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	err = Sd.Register(req.ServiceInstance)
	if err == nil {
		data, _, err = Sd.Discovery(context.TODO(), req.Service, "")
	}
}

type deregisterReq struct {
	Service string
	Id      string
}

func handleDeregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(deregisterReq)

	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(nil, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleDeregister: req:%+v error: %v", req, err)
			return
		}
		Sdlogger.Info("handleDeregister OK, req:%+v", req)
	}()

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()
	if req.Service == "" || req.Id == "" {
		code = 400
		err = errors.Wrap(ErrParams, "provide a valid service, id")
		return
	}
	err = Sd.Deregister(req.Service, req.Id)
	if err != nil {
		if errors.Is(err, ErrInstanceNotRegistered) {
			code = 400
		}
	}
}

type discoveryReq struct {
	Service   string
	LastHash  string
	WaitMaxMs int64
}
type discoveryRsp struct {
	Instances []ServiceInstance
	Hash      string
}

const MaxDiscoveryTimeout = time.Minute * 5

func handleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}

	body := new(discoveryReq)

	var err error
	var code = 200

	var response *discoveryRsp

	defer func() {
		rsp := ToJson(newRes(response, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleDiscovery: body:%+v error: %v", body, err)
			return
		}
		Sdlogger.Info("handleDiscovery OK, body:%+v", body)
	}()

	err = json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()
	if body.Service == "" {
		code = 400
		err = errors.Wrap(ErrParams, "need service")
		return
	}
	if body.WaitMaxMs > MaxDiscoveryTimeout.Milliseconds() {
		code = 400
		err = errors.Wrapf(ErrParams, "max wait ms is %s", MaxDiscoveryTimeout)
		return
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*time.Duration(body.WaitMaxMs))
	defer cancel()

	var (
		instances []ServiceInstance
		hash      string
	)
	instances, hash, err = Sd.Discovery(ctx, body.Service, body.LastHash)
	if err == nil {
		response = &discoveryRsp{
			Instances: instances,
			Hash:      hash,
		}
	}
}
