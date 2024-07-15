// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: github.com/dogmatiq/enginekit/protobuf/envelopepb/envelope.proto

package envelopepb

import (
	identitypb "github.com/dogmatiq/enginekit/protobuf/identitypb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Envelope is a container for a Dogma message and its meta-data.
type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// MessageId is a unique identifier for the message in this envelope.
	MessageId *uuidpb.UUID `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	// CausationId is the (optional) ID of the message that was the direct cause
	// of the message in this envelope.
	//
	// If it is the zero-value, the message was not caused by any other message.
	CausationId *uuidpb.UUID `protobuf:"bytes,2,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	// CorrelationId is the (optional) ID of the first ancestor of the message in
	// this envelope that was not caused by another message.
	//
	// If it is the zero-value, the message was not caused by any other message.
	CorrelationId *uuidpb.UUID `protobuf:"bytes,3,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	// SourceSite is the (optional) identity of the "site" that the source
	// application is running within.
	//
	// The site is used to disambiguate between messages from different
	// installations of the same application.
	SourceSite *identitypb.Identity `protobuf:"bytes,4,opt,name=source_site,json=sourceSite,proto3" json:"source_site,omitempty"`
	// SourceApplication is the identity of the Dogma application that produced
	// the message in this envelope.
	SourceApplication *identitypb.Identity `protobuf:"bytes,5,opt,name=source_application,json=sourceApplication,proto3" json:"source_application,omitempty"`
	// SourceHandler is the identity of the Dogma handler that produced the
	// message in this envelope.
	//
	// It is the zero-value if the message was not produced by a handler.
	SourceHandler *identitypb.Identity `protobuf:"bytes,6,opt,name=source_handler,json=sourceHandler,proto3" json:"source_handler,omitempty"`
	// SourceInstanceId is the ID of the aggregate or process instance that
	// produced the message in this envelope.
	//
	// It is empty if the message was not produced by an aggregate or process
	// handler.
	SourceInstanceId string `protobuf:"bytes,7,opt,name=source_instance_id,json=sourceInstanceId,proto3" json:"source_instance_id,omitempty"`
	// CreatedAt is the time at which the envelope was created.
	//
	// This is typically the point at which the message first enters the engine.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// ScheduledFor is the time at which a timeout message is scheduled to occur.
	//
	// It is the zero-value if the message is a command or event.
	ScheduledFor *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=scheduled_for,json=scheduledFor,proto3" json:"scheduled_for,omitempty"`
	// Description is a human-readable description of the message.
	Description string `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	// PortableName is the unique name used to identify messages of this type.
	PortableName string `protobuf:"bytes,11,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
	// MediaType is a MIME media-type describing the content and encoding of the
	// binary message data.
	MediaType string `protobuf:"bytes,12,opt,name=media_type,json=mediaType,proto3" json:"media_type,omitempty"`
	// Attributes is a set of arbitrary key/value pairs that provide additional
	// information about the message.
	//
	// Keys beginning with "_" are reserved for use by the enginekit module. All
	// other keys SHOULD use reverse-domain notation, e.g. "com.example.some-key".
	Attributes map[string]string `protobuf:"bytes,13,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Data is the binary message data.
	//
	// The data format is described by MediaType, the allowed values of both are
	// outside the scope of this specification.
	Data []byte `protobuf:"bytes,14,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescGZIP(), []int{0}
}

func (x *Envelope) GetMessageId() *uuidpb.UUID {
	if x != nil {
		return x.MessageId
	}
	return nil
}

func (x *Envelope) GetCausationId() *uuidpb.UUID {
	if x != nil {
		return x.CausationId
	}
	return nil
}

func (x *Envelope) GetCorrelationId() *uuidpb.UUID {
	if x != nil {
		return x.CorrelationId
	}
	return nil
}

func (x *Envelope) GetSourceSite() *identitypb.Identity {
	if x != nil {
		return x.SourceSite
	}
	return nil
}

func (x *Envelope) GetSourceApplication() *identitypb.Identity {
	if x != nil {
		return x.SourceApplication
	}
	return nil
}

func (x *Envelope) GetSourceHandler() *identitypb.Identity {
	if x != nil {
		return x.SourceHandler
	}
	return nil
}

func (x *Envelope) GetSourceInstanceId() string {
	if x != nil {
		return x.SourceInstanceId
	}
	return ""
}

func (x *Envelope) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Envelope) GetScheduledFor() *timestamppb.Timestamp {
	if x != nil {
		return x.ScheduledFor
	}
	return nil
}

func (x *Envelope) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Envelope) GetPortableName() string {
	if x != nil {
		return x.PortableName
	}
	return ""
}

func (x *Envelope) GetMediaType() string {
	if x != nil {
		return x.MediaType
	}
	return ""
}

func (x *Envelope) GetAttributes() map[string]string {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Envelope) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto protoreflect.FileDescriptor

var file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDesc = []byte{
	0x0a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x70, 0x62, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b,
	0x69, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x70, 0x62, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x75,
	0x69, 0x64, 0x70, 0x62, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xa7, 0x06, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x33, 0x0a, 0x0a,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x09, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49,
	0x64, 0x12, 0x37, 0x0a, 0x0c, 0x63, 0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x0b, 0x63,
	0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x3b, 0x0a, 0x0e, 0x63, 0x6f,
	0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x73, 0x69, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64,
	0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x69,
	0x74, 0x65, 0x12, 0x47, 0x0a, 0x12, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x70, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x11, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3f, 0x0a, 0x0e, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x0d, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x12,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3f, 0x0a, 0x0d, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x64, 0x5f, 0x66, 0x6f, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x64, 0x46, 0x6f, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x6f, 0x72, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x48, 0x0a, 0x0a,
	0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x28, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3d, 0x0a, 0x0f, 0x41, 0x74,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71,
	0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescData = file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDesc
)

func file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescData)
	})
	return file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_goTypes = []any{
	(*Envelope)(nil),              // 0: dogma.protobuf.Envelope
	nil,                           // 1: dogma.protobuf.Envelope.AttributesEntry
	(*uuidpb.UUID)(nil),           // 2: dogma.protobuf.UUID
	(*identitypb.Identity)(nil),   // 3: dogma.protobuf.Identity
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_depIdxs = []int32{
	2, // 0: dogma.protobuf.Envelope.message_id:type_name -> dogma.protobuf.UUID
	2, // 1: dogma.protobuf.Envelope.causation_id:type_name -> dogma.protobuf.UUID
	2, // 2: dogma.protobuf.Envelope.correlation_id:type_name -> dogma.protobuf.UUID
	3, // 3: dogma.protobuf.Envelope.source_site:type_name -> dogma.protobuf.Identity
	3, // 4: dogma.protobuf.Envelope.source_application:type_name -> dogma.protobuf.Identity
	3, // 5: dogma.protobuf.Envelope.source_handler:type_name -> dogma.protobuf.Identity
	4, // 6: dogma.protobuf.Envelope.created_at:type_name -> google.protobuf.Timestamp
	4, // 7: dogma.protobuf.Envelope.scheduled_for:type_name -> google.protobuf.Timestamp
	1, // 8: dogma.protobuf.Envelope.attributes:type_name -> dogma.protobuf.Envelope.AttributesEntry
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_init() }
func file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_init() {
	if File_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Envelope); i {
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
			RawDescriptor: file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto = out.File
	file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_rawDesc = nil
	file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_protobuf_envelopepb_envelope_proto_depIdxs = nil
}
