// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: flow.proto

package proto

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

type FlowEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique client event identifier
	EventId string `protobuf:"bytes,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	// Unique client flow session identifier
	FlowId string `protobuf:"bytes,2,opt,name=flow_id,json=flowId,proto3" json:"flow_id,omitempty"`
}

func (x *FlowEventRequest) Reset() {
	*x = FlowEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowEventRequest) ProtoMessage() {}

func (x *FlowEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowEventRequest.ProtoReflect.Descriptor instead.
func (*FlowEventRequest) Descriptor() ([]byte, []int) {
	return file_flow_proto_rawDescGZIP(), []int{0}
}

func (x *FlowEventRequest) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *FlowEventRequest) GetFlowId() string {
	if x != nil {
		return x.FlowId
	}
	return ""
}

type FlowEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique client event identifier that has been ack'ed
	EventId string `protobuf:"bytes,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
}

func (x *FlowEventResponse) Reset() {
	*x = FlowEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlowEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlowEventResponse) ProtoMessage() {}

func (x *FlowEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlowEventResponse.ProtoReflect.Descriptor instead.
func (*FlowEventResponse) Descriptor() ([]byte, []int) {
	return file_flow_proto_rawDescGZIP(), []int{1}
}

func (x *FlowEventResponse) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

var File_flow_proto protoreflect.FileDescriptor

var file_flow_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x46, 0x0a, 0x10, 0x46, 0x6c, 0x6f, 0x77,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x6c, 0x6f, 0x77, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6c, 0x6f, 0x77, 0x49, 0x64,
	0x22, 0x2e, 0x0a, 0x11, 0x46, 0x6c, 0x6f, 0x77, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x32, 0x5a, 0x0a, 0x0b, 0x46, 0x6c, 0x6f, 0x77, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4b, 0x0a, 0x06, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1c, 0x2e, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x08, 0x5a, 0x06,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flow_proto_rawDescOnce sync.Once
	file_flow_proto_rawDescData = file_flow_proto_rawDesc
)

func file_flow_proto_rawDescGZIP() []byte {
	file_flow_proto_rawDescOnce.Do(func() {
		file_flow_proto_rawDescData = protoimpl.X.CompressGZIP(file_flow_proto_rawDescData)
	})
	return file_flow_proto_rawDescData
}

var file_flow_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_flow_proto_goTypes = []interface{}{
	(*FlowEventRequest)(nil),  // 0: management.FlowEventRequest
	(*FlowEventResponse)(nil), // 1: management.FlowEventResponse
}
var file_flow_proto_depIdxs = []int32{
	0, // 0: management.FlowService.Events:input_type -> management.FlowEventRequest
	1, // 1: management.FlowService.Events:output_type -> management.FlowEventResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_flow_proto_init() }
func file_flow_proto_init() {
	if File_flow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_flow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowEventRequest); i {
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
		file_flow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlowEventResponse); i {
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
			RawDescriptor: file_flow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_flow_proto_goTypes,
		DependencyIndexes: file_flow_proto_depIdxs,
		MessageInfos:      file_flow_proto_msgTypes,
	}.Build()
	File_flow_proto = out.File
	file_flow_proto_rawDesc = nil
	file_flow_proto_goTypes = nil
	file_flow_proto_depIdxs = nil
}
