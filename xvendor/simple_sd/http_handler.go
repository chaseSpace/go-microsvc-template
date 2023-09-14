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

type PingReq struct {
	Ping bool
}
type PingRspBody struct {
	Pong bool
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}

	req := new(PingReq)

	var err error
	var code = 200

	var response *PingRspBody

	defer func() {
		rsp := ToJson(newRes(response, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handlePing: req:%+v error: %v", req, err)
			return
		}
		Sdlogger.Debug("handlePing OK, req:%+v", req)
	}()

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		code = 400
		err = errors.Wrap(ErrParams, err.Error())
		return
	}
	_ = r.Body.Close()
	if req.Ping {
		response = &PingRspBody{Pong: true}
	}
}

type RegisterReq struct {
	ServiceInstance
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(RegisterReq)

	st := time.Now()

	var data []ServiceInstance
	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(data, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleRegister: service:%s req:%s error: %v", req.Name, req.Addr(), err)
			return
		}
		Sdlogger.Debug("handleRegister OK, dur:%s service:%s req:%s", time.Since(st), req.Name, req.Addr())
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
		data, _, err = Sd.Discovery(context.TODO(), req.Name, "")
	}
}

type DeregisterReq struct {
	Service string
	Id      string
}

func handleDeregister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}
	req := new(DeregisterReq)

	var err error
	var code = 200
	defer func() {
		rsp := ToJson(newRes(nil, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleDeregister: req:%+v error: %v", req, err)
			return
		}
		Sdlogger.Debug("handleDeregister OK, req:%+v", req)
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

type DiscoveryReq struct {
	Service   string
	LastHash  string
	WaitMaxMs int64
}
type DiscoveryRspBody struct {
	Instances []ServiceInstance
	Hash      string
}

const MaxDiscoveryTimeout = time.Minute * 5

func handleDiscovery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}

	body := new(DiscoveryReq)

	var err error
	var code = 200

	var response *DiscoveryRspBody

	st := time.Now()
	defer func() {
		rsp := ToJson(newRes(response, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleDiscovery: body:%+v error: %v", body, err)
			return
		}
		Sdlogger.Debug("handleDiscovery OK, dur:%dms  body:%+v  --rsp:%s", time.Since(st).Milliseconds(), body, rsp)
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
		err = errors.Wrap(ErrParams, "params without service")
		return
	}
	if body.WaitMaxMs > MaxDiscoveryTimeout.Milliseconds() {
		code = 400
		err = errors.Wrapf(ErrParams, "max wait-time ms is %d", MaxDiscoveryTimeout.Milliseconds())
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
		response = &DiscoveryRspBody{
			Instances: instances,
			Hash:      hash,
		}
	}
}

type HealthCheckReq struct {
	Service string
	Id      string
}
type HealthCheckRspBody struct {
	Registered bool
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		_, _ = w.Write(ToJson(newRes(nil, 400, ErrMethod)))
		return
	}

	req := new(HealthCheckReq)
	rspBody := new(HealthCheckRspBody)

	var err error
	var code = 200

	st := time.Now()
	defer func() {
		rsp := ToJson(newRes(rspBody, code, err))
		_, _ = w.Write(rsp)
		if err != nil {
			Sdlogger.Error("handleHealthCheck: req:%+v error: %v", req, err)
			return
		}
		Sdlogger.Debug("handleHealthCheck OK, dur:%dms  req:%+v  --rsp:%s", time.Since(st).Milliseconds(), req, rsp)
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
		err = errors.Wrap(ErrParams, "provide a valid service and id")
		return
	}

	rspBody.Registered = Sd.HealthCheck(req.Service, req.Id)
}
