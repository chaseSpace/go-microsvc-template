package protobytes

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

const Name = "bytes"

func init() {
	encoding.RegisterCodec(codecBytes{})
}

type codecBytes struct{}

func (codecBytes) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		vb, ok := v.([]byte) // from gateway
		if !ok {
			return nil, fmt.Errorf("failed to marshal, message is %T, want protobytes.Message or []byte", v)
		}
		return vb, nil
	}
	return jsoniter.Marshal(vv)
}

func (codecBytes) Unmarshal(data []byte, v interface{}) error {
	vv, ok := v.(proto.Message)
	if !ok {
		vb, ok := v.(*bytes.Buffer) // from gateway
		if !ok {
			return fmt.Errorf("failed to marshal, message is %T, want protobytes.Message or []byte", v)
		}
		_, err := vb.Write(data)
		return err
	}
	return jsoniter.Unmarshal(data, vv)
}

func (codecBytes) Name() string {
	return Name
}
