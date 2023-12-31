// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.4
// source: svc/user/user.ext.proto

package user

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	svc "microsvc/protocol/svc"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetUserReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base *svc.BaseExtReq `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"` // 每个外部接口Req都必须添加这个成员类型，grpc拦截器会做验证
	Uids []int64         `protobuf:"varint,2,rep,packed,name=uids,proto3" json:"uids,omitempty"`
}

func (x *GetUserReq) Reset() {
	*x = GetUserReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserReq) ProtoMessage() {}

func (x *GetUserReq) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserReq.ProtoReflect.Descriptor instead.
func (*GetUserReq) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{0}
}

func (x *GetUserReq) GetBase() *svc.BaseExtReq {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *GetUserReq) GetUids() []int64 {
	if x != nil {
		return x.Uids
	}
	return nil
}

type GetUserRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Umap map[int64]*User `protobuf:"bytes,1,rep,name=umap,proto3" json:"umap,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GetUserRes) Reset() {
	*x = GetUserRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRes) ProtoMessage() {}

func (x *GetUserRes) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRes.ProtoReflect.Descriptor instead.
func (*GetUserRes) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserRes) GetUmap() map[int64]*User {
	if x != nil {
		return x.Umap
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid      int64  `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Nickname string `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Birthday string `protobuf:"bytes,3,opt,name=birthday,proto3" json:"birthday,omitempty"`
	Sex      int32  `protobuf:"varint,4,opt,name=sex,proto3" json:"sex,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{2}
}

func (x *User) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *User) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *User) GetBirthday() string {
	if x != nil {
		return x.Birthday
	}
	return ""
}

func (x *User) GetSex() int32 {
	if x != nil {
		return x.Sex
	}
	return 0
}

type SignUpReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nickname      string `protobuf:"bytes,1,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Sex           int32  `protobuf:"varint,2,opt,name=sex,proto3" json:"sex,omitempty"`
	Birthday      string `protobuf:"bytes,3,opt,name=birthday,proto3" json:"birthday,omitempty"`
	PhoneAreaCode string `protobuf:"bytes,4,opt,name=phone_area_code,json=phoneAreaCode,proto3" json:"phone_area_code,omitempty"`
	Phone         string `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	VerifyCode    string `protobuf:"bytes,6,opt,name=verify_code,json=verifyCode,proto3" json:"verify_code,omitempty"`
}

func (x *SignUpReq) Reset() {
	*x = SignUpReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignUpReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignUpReq) ProtoMessage() {}

func (x *SignUpReq) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignUpReq.ProtoReflect.Descriptor instead.
func (*SignUpReq) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{3}
}

func (x *SignUpReq) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *SignUpReq) GetSex() int32 {
	if x != nil {
		return x.Sex
	}
	return 0
}

func (x *SignUpReq) GetBirthday() string {
	if x != nil {
		return x.Birthday
	}
	return ""
}

func (x *SignUpReq) GetPhoneAreaCode() string {
	if x != nil {
		return x.PhoneAreaCode
	}
	return ""
}

func (x *SignUpReq) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *SignUpReq) GetVerifyCode() string {
	if x != nil {
		return x.VerifyCode
	}
	return ""
}

type SignUpRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *SignUpRes) Reset() {
	*x = SignUpRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignUpRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignUpRes) ProtoMessage() {}

func (x *SignUpRes) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignUpRes.ProtoReflect.Descriptor instead.
func (*SignUpRes) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{4}
}

func (x *SignUpRes) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type SignInReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base          *svc.BaseExtReq `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	PhoneAreaCode string          `protobuf:"bytes,2,opt,name=phone_area_code,json=phoneAreaCode,proto3" json:"phone_area_code,omitempty"`
	Phone         string          `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	VerifyCode    string          `protobuf:"bytes,4,opt,name=verify_code,json=verifyCode,proto3" json:"verify_code,omitempty"`
}

func (x *SignInReq) Reset() {
	*x = SignInReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInReq) ProtoMessage() {}

func (x *SignInReq) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInReq.ProtoReflect.Descriptor instead.
func (*SignInReq) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{5}
}

func (x *SignInReq) GetBase() *svc.BaseExtReq {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *SignInReq) GetPhoneAreaCode() string {
	if x != nil {
		return x.PhoneAreaCode
	}
	return ""
}

func (x *SignInReq) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *SignInReq) GetVerifyCode() string {
	if x != nil {
		return x.VerifyCode
	}
	return ""
}

type SignInRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info  *User  `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *SignInRes) Reset() {
	*x = SignInRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_user_user_ext_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignInRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignInRes) ProtoMessage() {}

func (x *SignInRes) ProtoReflect() protoreflect.Message {
	mi := &file_svc_user_user_ext_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignInRes.ProtoReflect.Descriptor instead.
func (*SignInRes) Descriptor() ([]byte, []int) {
	return file_svc_user_user_ext_proto_rawDescGZIP(), []int{6}
}

func (x *SignInRes) GetInfo() *User {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *SignInRes) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_svc_user_user_ext_proto protoreflect.FileDescriptor

var file_svc_user_user_ext_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x76, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x76, 0x63, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x1a, 0x10, 0x73, 0x76, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x12, 0x23, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x76, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x45, 0x78, 0x74, 0x52,
	0x65, 0x71, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x69, 0x64, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x04, 0x75, 0x69, 0x64, 0x73, 0x22, 0x89, 0x01, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x04, 0x75,
	0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x76, 0x63, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x2e,
	0x55, 0x6d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x75, 0x6d, 0x61, 0x70, 0x1a,
	0x47, 0x0a, 0x09, 0x55, 0x6d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x24,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x73, 0x76, 0x63, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x62, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75,
	0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64, 0x61, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x62, 0x69, 0x72, 0x74, 0x68, 0x64, 0x61, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65,
	0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x73, 0x65, 0x78, 0x22, 0xb4, 0x01, 0x0a,
	0x09, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69,
	0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69,
	0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x78, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x73, 0x65, 0x78, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x69, 0x72, 0x74,
	0x68, 0x64, 0x61, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x69, 0x72, 0x74,
	0x68, 0x64, 0x61, 0x79, 0x12, 0x26, 0x0a, 0x0f, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f, 0x61, 0x72,
	0x65, 0x61, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70,
	0x68, 0x6f, 0x6e, 0x65, 0x41, 0x72, 0x65, 0x61, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x43,
	0x6f, 0x64, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x8f, 0x01, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x49,
	0x6e, 0x52, 0x65, 0x71, 0x12, 0x23, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x76, 0x63, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x45, 0x78, 0x74,
	0x52, 0x65, 0x71, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x5f, 0x61, 0x72, 0x65, 0x61, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x41, 0x72, 0x65, 0x61, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x45, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e,
	0x49, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x76, 0x63, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x32,
	0xa8, 0x01, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x45, 0x78, 0x74, 0x12, 0x32, 0x0a, 0x06, 0x53,
	0x69, 0x67, 0x6e, 0x55, 0x70, 0x12, 0x13, 0x2e, 0x73, 0x76, 0x63, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x73, 0x76, 0x63,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x70, 0x52, 0x65, 0x73, 0x12,
	0x32, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x12, 0x13, 0x2e, 0x73, 0x76, 0x63, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x13,
	0x2e, 0x73, 0x76, 0x63, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e,
	0x52, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14,
	0x2e, 0x73, 0x76, 0x63, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x73, 0x76, 0x63, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x42, 0x1c, 0x5a, 0x1a, 0x6d, 0x69,
	0x63, 0x72, 0x6f, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f,
	0x73, 0x76, 0x63, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_svc_user_user_ext_proto_rawDescOnce sync.Once
	file_svc_user_user_ext_proto_rawDescData = file_svc_user_user_ext_proto_rawDesc
)

func file_svc_user_user_ext_proto_rawDescGZIP() []byte {
	file_svc_user_user_ext_proto_rawDescOnce.Do(func() {
		file_svc_user_user_ext_proto_rawDescData = protoimpl.X.CompressGZIP(file_svc_user_user_ext_proto_rawDescData)
	})
	return file_svc_user_user_ext_proto_rawDescData
}

var file_svc_user_user_ext_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_svc_user_user_ext_proto_goTypes = []interface{}{
	(*GetUserReq)(nil),     // 0: svc.user.GetUserReq
	(*GetUserRes)(nil),     // 1: svc.user.GetUserRes
	(*User)(nil),           // 2: svc.user.User
	(*SignUpReq)(nil),      // 3: svc.user.SignUpReq
	(*SignUpRes)(nil),      // 4: svc.user.SignUpRes
	(*SignInReq)(nil),      // 5: svc.user.SignInReq
	(*SignInRes)(nil),      // 6: svc.user.SignInRes
	nil,                    // 7: svc.user.GetUserRes.UmapEntry
	(*svc.BaseExtReq)(nil), // 8: svc.BaseExtReq
}
var file_svc_user_user_ext_proto_depIdxs = []int32{
	8, // 0: svc.user.GetUserReq.base:type_name -> svc.BaseExtReq
	7, // 1: svc.user.GetUserRes.umap:type_name -> svc.user.GetUserRes.UmapEntry
	8, // 2: svc.user.SignInReq.base:type_name -> svc.BaseExtReq
	2, // 3: svc.user.SignInRes.info:type_name -> svc.user.User
	2, // 4: svc.user.GetUserRes.UmapEntry.value:type_name -> svc.user.User
	3, // 5: svc.user.UserExt.SignUp:input_type -> svc.user.SignUpReq
	5, // 6: svc.user.UserExt.SignIn:input_type -> svc.user.SignInReq
	0, // 7: svc.user.UserExt.GetUser:input_type -> svc.user.GetUserReq
	4, // 8: svc.user.UserExt.SignUp:output_type -> svc.user.SignUpRes
	6, // 9: svc.user.UserExt.SignIn:output_type -> svc.user.SignInRes
	1, // 10: svc.user.UserExt.GetUser:output_type -> svc.user.GetUserRes
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_svc_user_user_ext_proto_init() }
func file_svc_user_user_ext_proto_init() {
	if File_svc_user_user_ext_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_svc_user_user_ext_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserReq); i {
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
		file_svc_user_user_ext_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserRes); i {
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
		file_svc_user_user_ext_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_svc_user_user_ext_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignUpReq); i {
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
		file_svc_user_user_ext_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignUpRes); i {
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
		file_svc_user_user_ext_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInReq); i {
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
		file_svc_user_user_ext_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignInRes); i {
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
			RawDescriptor: file_svc_user_user_ext_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_svc_user_user_ext_proto_goTypes,
		DependencyIndexes: file_svc_user_user_ext_proto_depIdxs,
		MessageInfos:      file_svc_user_user_ext_proto_msgTypes,
	}.Build()
	File_svc_user_user_ext_proto = out.File
	file_svc_user_user_ext_proto_rawDesc = nil
	file_svc_user_user_ext_proto_goTypes = nil
	file_svc_user_user_ext_proto_depIdxs = nil
}
