// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.0
// source: svc/common.proto

package svc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BaseExtReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	App        string `protobuf:"bytes,1,opt,name=app,proto3" json:"app,omitempty"`
	AppVersion string `protobuf:"bytes,2,opt,name=app_version,json=appVersion,proto3" json:"app_version,omitempty"` //...
}

func (x *BaseExtReq) Reset() {
	*x = BaseExtReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseExtReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseExtReq) ProtoMessage() {}

func (x *BaseExtReq) ProtoReflect() protoreflect.Message {
	mi := &file_svc_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseExtReq.ProtoReflect.Descriptor instead.
func (*BaseExtReq) Descriptor() ([]byte, []int) {
	return file_svc_common_proto_rawDescGZIP(), []int{0}
}

func (x *BaseExtReq) GetApp() string {
	if x != nil {
		return x.App
	}
	return ""
}

func (x *BaseExtReq) GetAppVersion() string {
	if x != nil {
		return x.AppVersion
	}
	return ""
}

type AdminBaseReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid  int64  `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Nick string `protobuf:"bytes,2,opt,name=nick,proto3" json:"nick,omitempty"` // ...
}

func (x *AdminBaseReq) Reset() {
	*x = AdminBaseReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AdminBaseReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminBaseReq) ProtoMessage() {}

func (x *AdminBaseReq) ProtoReflect() protoreflect.Message {
	mi := &file_svc_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminBaseReq.ProtoReflect.Descriptor instead.
func (*AdminBaseReq) Descriptor() ([]byte, []int) {
	return file_svc_common_proto_rawDescGZIP(), []int{1}
}

func (x *AdminBaseReq) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *AdminBaseReq) GetNick() string {
	if x != nil {
		return x.Nick
	}
	return ""
}

type AdminCommonRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32      `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string     `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *anypb.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *AdminCommonRsp) Reset() {
	*x = AdminCommonRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AdminCommonRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminCommonRsp) ProtoMessage() {}

func (x *AdminCommonRsp) ProtoReflect() protoreflect.Message {
	mi := &file_svc_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminCommonRsp.ProtoReflect.Descriptor instead.
func (*AdminCommonRsp) Descriptor() ([]byte, []int) {
	return file_svc_common_proto_rawDescGZIP(), []int{2}
}

func (x *AdminCommonRsp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *AdminCommonRsp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *AdminCommonRsp) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

type HttpCommonRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32      `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string     `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *anypb.Any `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *HttpCommonRsp) Reset() {
	*x = HttpCommonRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpCommonRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpCommonRsp) ProtoMessage() {}

func (x *HttpCommonRsp) ProtoReflect() protoreflect.Message {
	mi := &file_svc_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpCommonRsp.ProtoReflect.Descriptor instead.
func (*HttpCommonRsp) Descriptor() ([]byte, []int) {
	return file_svc_common_proto_rawDescGZIP(), []int{3}
}

func (x *HttpCommonRsp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *HttpCommonRsp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *HttpCommonRsp) GetData() *anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_svc_common_proto protoreflect.FileDescriptor

var file_svc_common_proto_rawDesc = []byte{
	0x0a, 0x10, 0x73, 0x76, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x73, 0x76, 0x63, 0x1a, 0x21, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65,
	0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x0a, 0x42, 0x61,
	0x73, 0x65, 0x45, 0x78, 0x74, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x70, 0x70, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x70,
	0x70, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x61, 0x70, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x34, 0x0a, 0x0c, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x42, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x69, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x69, 0x63,
	0x6b, 0x22, 0x60, 0x0a, 0x0e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x52, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x5f, 0x0a, 0x0d, 0x48, 0x74, 0x74, 0x70, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x52, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x42, 0x17, 0x5a, 0x15, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x76, 0x63,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x73, 0x76, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_svc_common_proto_rawDescOnce sync.Once
	file_svc_common_proto_rawDescData = file_svc_common_proto_rawDesc
)

func file_svc_common_proto_rawDescGZIP() []byte {
	file_svc_common_proto_rawDescOnce.Do(func() {
		file_svc_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_svc_common_proto_rawDescData)
	})
	return file_svc_common_proto_rawDescData
}

var file_svc_common_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_svc_common_proto_goTypes = []interface{}{
	(*BaseExtReq)(nil),     // 0: svc.BaseExtReq
	(*AdminBaseReq)(nil),   // 1: svc.AdminBaseReq
	(*AdminCommonRsp)(nil), // 2: svc.AdminCommonRsp
	(*HttpCommonRsp)(nil),  // 3: svc.HttpCommonRsp
	(*anypb.Any)(nil),      // 4: google.protobuf.Any
}
var file_svc_common_proto_depIdxs = []int32{
	4, // 0: svc.AdminCommonRsp.data:type_name -> google.protobuf.Any
	4, // 1: svc.HttpCommonRsp.data:type_name -> google.protobuf.Any
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_svc_common_proto_init() }
func file_svc_common_proto_init() {
	if File_svc_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_svc_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseExtReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_svc_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AdminBaseReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_svc_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AdminCommonRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_svc_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpCommonRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_svc_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_svc_common_proto_goTypes,
		DependencyIndexes: file_svc_common_proto_depIdxs,
		MessageInfos:      file_svc_common_proto_msgTypes,
	}.Build()
	File_svc_common_proto = out.File
	file_svc_common_proto_rawDesc = nil
	file_svc_common_proto_goTypes = nil
	file_svc_common_proto_depIdxs = nil
}
