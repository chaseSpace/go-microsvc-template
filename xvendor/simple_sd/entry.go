package simple_sd

import (
	"fmt"
	"net/http"
)

type HTTPServer struct {
}

func NewSimpleSdHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Run(port int) error {
	http.HandleFunc("/service/register", handleRegister)
	http.HandleFunc("/service/deregister", handleDeregister)
	http.HandleFunc("/service/discovery", handleDiscovery)

	Sdlogger.Info("SimpleSdHTTPServer is running on http://localhost:%d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}
