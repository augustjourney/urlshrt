// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: urls.proto

package __

import (
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

type CreateURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalUrl string `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	UserId      string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *CreateURLRequest) Reset() {
	*x = CreateURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_urls_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateURLRequest) ProtoMessage() {}

func (x *CreateURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urls_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateURLRequest.ProtoReflect.Descriptor instead.
func (*CreateURLRequest) Descriptor() ([]byte, []int) {
	return file_urls_proto_rawDescGZIP(), []int{0}
}

func (x *CreateURLRequest) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *CreateURLRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type CreateURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortUrl string `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
}

func (x *CreateURLResponse) Reset() {
	*x = CreateURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_urls_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateURLResponse) ProtoMessage() {}

func (x *CreateURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urls_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateURLResponse.ProtoReflect.Descriptor instead.
func (*CreateURLResponse) Descriptor() ([]byte, []int) {
	return file_urls_proto_rawDescGZIP(), []int{1}
}

func (x *CreateURLResponse) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type GetURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortUrl string `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
}

func (x *GetURLRequest) Reset() {
	*x = GetURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_urls_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLRequest) ProtoMessage() {}

func (x *GetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_urls_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLRequest.ProtoReflect.Descriptor instead.
func (*GetURLRequest) Descriptor() ([]byte, []int) {
	return file_urls_proto_rawDescGZIP(), []int{2}
}

func (x *GetURLRequest) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type GetURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalUrl string `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *GetURLResponse) Reset() {
	*x = GetURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_urls_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLResponse) ProtoMessage() {}

func (x *GetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_urls_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLResponse.ProtoReflect.Descriptor instead.
func (*GetURLResponse) Descriptor() ([]byte, []int) {
	return file_urls_proto_rawDescGZIP(), []int{3}
}

func (x *GetURLResponse) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

var File_urls_proto protoreflect.FileDescriptor

var file_urls_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x75, 0x72, 0x6c, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4e, 0x0a, 0x10,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c,
	0x55, 0x72, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x30, 0x0a, 0x11,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x22, 0x2c,
	0x0a, 0x0d, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x22, 0x33, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21,
	0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72,
	0x6c, 0x32, 0x64, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x12, 0x32, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x55, 0x52, 0x4c, 0x12, 0x11, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06,
	0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x0e, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_urls_proto_rawDescOnce sync.Once
	file_urls_proto_rawDescData = file_urls_proto_rawDesc
)

func file_urls_proto_rawDescGZIP() []byte {
	file_urls_proto_rawDescOnce.Do(func() {
		file_urls_proto_rawDescData = protoimpl.X.CompressGZIP(file_urls_proto_rawDescData)
	})
	return file_urls_proto_rawDescData
}

var file_urls_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_urls_proto_goTypes = []any{
	(*CreateURLRequest)(nil),  // 0: CreateURLRequest
	(*CreateURLResponse)(nil), // 1: CreateURLResponse
	(*GetURLRequest)(nil),     // 2: GetURLRequest
	(*GetURLResponse)(nil),    // 3: GetURLResponse
}
var file_urls_proto_depIdxs = []int32{
	0, // 0: URL.CreateURL:input_type -> CreateURLRequest
	2, // 1: URL.GetURL:input_type -> GetURLRequest
	1, // 2: URL.CreateURL:output_type -> CreateURLResponse
	3, // 3: URL.GetURL:output_type -> GetURLResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_urls_proto_init() }
func file_urls_proto_init() {
	if File_urls_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_urls_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateURLRequest); i {
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
		file_urls_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateURLResponse); i {
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
		file_urls_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GetURLRequest); i {
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
		file_urls_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetURLResponse); i {
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
			RawDescriptor: file_urls_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_urls_proto_goTypes,
		DependencyIndexes: file_urls_proto_depIdxs,
		MessageInfos:      file_urls_proto_msgTypes,
	}.Build()
	File_urls_proto = out.File
	file_urls_proto_rawDesc = nil
	file_urls_proto_goTypes = nil
	file_urls_proto_depIdxs = nil
}