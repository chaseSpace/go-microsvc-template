package xhttp

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"microsvc/pkg/xlog"
)

type FastHttpSvr struct {
	port    int
	handler fasthttp.RequestHandler
}

func New(port int, handler fasthttp.RequestHandler) *FastHttpSvr {
	return &FastHttpSvr{port: port, handler: handler}
}

func (f FastHttpSvr) Start() {
	fmt.Println("\nCongratulations! ^_^")
	fmt.Printf("serving HTTP on http://localhost:%v\n\n", f.port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", f.port), f.handler)
	if err != nil {
		xlog.Error(fmt.Sprintf("xhttp: failed to server HTTP on: localhost:%d", f.port))
	}
}
