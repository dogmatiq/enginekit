// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.1
// source: github.com/dogmatiq/enginekit/protobuf/identitypb/identity.proto

package identitypb

import (
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
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

// Identity represents the identity of an entity.
//
// It is used to identify Dogma applications, handlers, sites, streams, etc.
type Identity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name is the entity's human-readable name.
	//
	// The name should be unique enough to allow a human to identify the entity
	// without ambiguity. There is no hard requirement for uniqueness.
	//
	// Entity names may be changed at any time.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Key is the entity's immutable, unique key.
	//
	// The key is used to uniquely identify the entity globally, and across all
	// time. Every entity must have its own UUID, which must not be changed.
	Key *uuidpb.UUID `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *Identity) Reset() {
	*x = Identity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Identity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Identity) ProtoMessage() {}

func (x *Identity) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Identity.ProtoReflect.Descriptor instead.
func (*Identity) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescGZIP(), []int{0}
}

func (x *Identity) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Identity) GetKey() *uuidpb.UUID {
	if x != nil {
		return x.Key
	}
	return nil
}

var File_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto protoreflect.FileDescriptor

var file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDesc = []byte{
	0x0a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x70, 0x62, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x1a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64,
	0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69,
	0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x70,
	0x62, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x08,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d,
	0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescData = file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDesc
)

func file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescData)
	})
	return file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_goTypes = []interface{}{
	(*Identity)(nil),    // 0: dogma.protobuf.Identity
	(*uuidpb.UUID)(nil), // 1: dogma.protobuf.UUID
}
var file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_depIdxs = []int32{
	1, // 0: dogma.protobuf.Identity.key:type_name -> dogma.protobuf.UUID
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_init() }
func file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_init() {
	if File_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Identity); i {
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
			RawDescriptor: file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto = out.File
	file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_rawDesc = nil
	file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_protobuf_identitypb_identity_proto_depIdxs = nil
}
