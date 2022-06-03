// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: discogs.proto

package discogs

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestTokenParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RequestTokenParams) Reset() {
	*x = RequestTokenParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_discogs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestTokenParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestTokenParams) ProtoMessage() {}

func (x *RequestTokenParams) ProtoReflect() protoreflect.Message {
	mi := &file_discogs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestTokenParams.ProtoReflect.Descriptor instead.
func (*RequestTokenParams) Descriptor() ([]byte, []int) {
	return file_discogs_proto_rawDescGZIP(), []int{0}
}

type RequestTokenResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *RequestTokenResult) Reset() {
	*x = RequestTokenResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_discogs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestTokenResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestTokenResult) ProtoMessage() {}

func (x *RequestTokenResult) ProtoReflect() protoreflect.Message {
	mi := &file_discogs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestTokenResult.ProtoReflect.Descriptor instead.
func (*RequestTokenResult) Descriptor() ([]byte, []int) {
	return file_discogs_proto_rawDescGZIP(), []int{1}
}

func (x *RequestTokenResult) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AccessTokenParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestToken string `protobuf:"bytes,1,opt,name=requestToken,proto3" json:"requestToken,omitempty"`
}

func (x *AccessTokenParams) Reset() {
	*x = AccessTokenParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_discogs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessTokenParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessTokenParams) ProtoMessage() {}

func (x *AccessTokenParams) ProtoReflect() protoreflect.Message {
	mi := &file_discogs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessTokenParams.ProtoReflect.Descriptor instead.
func (*AccessTokenParams) Descriptor() ([]byte, []int) {
	return file_discogs_proto_rawDescGZIP(), []int{2}
}

func (x *AccessTokenParams) GetRequestToken() string {
	if x != nil {
		return x.RequestToken
	}
	return ""
}

type AccessTokenResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token  string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Secret string `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (x *AccessTokenResult) Reset() {
	*x = AccessTokenResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_discogs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessTokenResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessTokenResult) ProtoMessage() {}

func (x *AccessTokenResult) ProtoReflect() protoreflect.Message {
	mi := &file_discogs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessTokenResult.ProtoReflect.Descriptor instead.
func (*AccessTokenResult) Descriptor() ([]byte, []int) {
	return file_discogs_proto_rawDescGZIP(), []int{3}
}

func (x *AccessTokenResult) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *AccessTokenResult) GetSecret() string {
	if x != nil {
		return x.Secret
	}
	return ""
}

var File_discogs_proto protoreflect.FileDescriptor

var file_discogs_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67,
	0x65, 0x6e, 0x2d, 0x73, 0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x14, 0x0a, 0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x22, 0x2a, 0x0a, 0x12, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x37, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x41, 0x0a, 0x11, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x32, 0xf3, 0x02, 0x0a, 0x07, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x12, 0xba,
	0x01, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1b, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x1b, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x70, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x12, 0x12, 0x10, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x67, 0x73, 0x92, 0x41, 0x55, 0x12, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x20,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x1a, 0x44, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x20, 0x74,
	0x68, 0x65, 0x20, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x20, 0x75, 0x73, 0x65, 0x64, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x74, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x68, 0x65, 0x20, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73,
	0x20, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x20, 0x66, 0x6c, 0x6f, 0x77, 0x12, 0xaa, 0x01, 0x0a, 0x0b,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x2e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x67, 0x73, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x1a, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x67,
	0x73, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x63, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x22, 0x10, 0x2f, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x92, 0x41, 0x48,
	0x12, 0x0c, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x20, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x1a, 0x38,
	0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x20, 0x74, 0x68, 0x65, 0x20, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x20, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x20, 0x75, 0x73, 0x65, 0x64, 0x20, 0x66, 0x6f,
	0x72, 0x20, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x20, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x69, 0x6e, 0x6e, 0x65, 0x6d, 0x61, 0x6e, 0x61,
	0x2f, 0x76, 0x69, 0x6e, 0x79, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x67, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_discogs_proto_rawDescOnce sync.Once
	file_discogs_proto_rawDescData = file_discogs_proto_rawDesc
)

func file_discogs_proto_rawDescGZIP() []byte {
	file_discogs_proto_rawDescOnce.Do(func() {
		file_discogs_proto_rawDescData = protoimpl.X.CompressGZIP(file_discogs_proto_rawDescData)
	})
	return file_discogs_proto_rawDescData
}

var file_discogs_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_discogs_proto_goTypes = []interface{}{
	(*RequestTokenParams)(nil), // 0: discogs.RequestTokenParams
	(*RequestTokenResult)(nil), // 1: discogs.RequestTokenResult
	(*AccessTokenParams)(nil),  // 2: discogs.AccessTokenParams
	(*AccessTokenResult)(nil),  // 3: discogs.AccessTokenResult
}
var file_discogs_proto_depIdxs = []int32{
	0, // 0: discogs.Discogs.RequestToken:input_type -> discogs.RequestTokenParams
	2, // 1: discogs.Discogs.AccessToken:input_type -> discogs.AccessTokenParams
	1, // 2: discogs.Discogs.RequestToken:output_type -> discogs.RequestTokenResult
	3, // 3: discogs.Discogs.AccessToken:output_type -> discogs.AccessTokenResult
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_discogs_proto_init() }
func file_discogs_proto_init() {
	if File_discogs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_discogs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestTokenParams); i {
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
		file_discogs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestTokenResult); i {
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
		file_discogs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessTokenParams); i {
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
		file_discogs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessTokenResult); i {
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
			RawDescriptor: file_discogs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_discogs_proto_goTypes,
		DependencyIndexes: file_discogs_proto_depIdxs,
		MessageInfos:      file_discogs_proto_msgTypes,
	}.Build()
	File_discogs_proto = out.File
	file_discogs_proto_rawDesc = nil
	file_discogs_proto_goTypes = nil
	file_discogs_proto_depIdxs = nil
}