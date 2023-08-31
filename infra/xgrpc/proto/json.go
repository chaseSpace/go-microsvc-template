package proto

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

// Json is the name registered for the proto compressor.
const Json = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message", v)
	}
	return json.Marshal(vv)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}
	return json.Unmarshal(data, vv)
}

func (codec) Name() string {
	return Json
}
