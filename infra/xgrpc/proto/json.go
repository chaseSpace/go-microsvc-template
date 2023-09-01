package proto

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

// Json is the name registered for the proto compressor.
const Json = "json"

func init() {
	encoding.RegisterCodec(codec{})
}

type ArbitraryBody map[string]interface{}

// codec is a Codec implementation with protobuf. It is the default codec for gRPC.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	_, ok := v.(proto.Message)
	if !ok {
		_, ok = v.(*ArbitraryBody) // from gateway
		if !ok {
			return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message or ArbitraryBody", v)
		}
	}
	return json.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	_, ok := v.(proto.Message)
	if !ok {
		_, ok = v.(*ArbitraryBody) // from gateway
		if !ok {
			return fmt.Errorf("failed to Unmarshal, message is %T, want proto.Message or ArbitraryBody", v)
		}
	}
	return json.Unmarshal(data, v)
}

func (codec) Name() string {
	return Json
}
