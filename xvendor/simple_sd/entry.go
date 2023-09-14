package simple_sd

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
)

type HTTPServer struct {
	port int
}

func NewSimpleSdHTTPServer(port int) *HTTPServer {
	return &HTTPServer{port: port}
}

func (h *HTTPServer) IsRunningOnLocal() bool {
	body := new(PingRspBody)
	rsp := newRes(body, 0, nil)
	_, _, errs := gorequest.New().Post(fmt.Sprintf("http://localhost:%d/ping", h.port)).SendStruct(&PingReq{Ping: true}).EndStruct(rsp)
	if len(errs) == 0 {
		return body.Pong
	}
	return false
}

func (h *HTTPServer) Run() error {
	Init()

	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/service/register", handleRegister)
	http.HandleFunc("/service/deregister", handleDeregister)
	http.HandleFunc("/service/discovery", handleDiscovery)
	http.HandleFunc("/service/health_check", handleHealthCheck)

	//Sdlogger.Info("SimpleSdHTTPServer is running on http://localhost:%d", h.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", h.port), nil)

}
