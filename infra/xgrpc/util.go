package xgrpc

import "bytes"

func beautifyReqAndResInClient(req, reply interface{}) (interface{}, interface{}) {
	// When grpc call is from gateway,the request is []byte, and the reply is *bytes.Buffer
	_req, _ := req.([]byte)
	if _req != nil {
		req = string(_req)
	}
	res, _ := reply.(*bytes.Buffer)
	if res != nil {
		reply = res.String()
	}
	return req, reply
}
