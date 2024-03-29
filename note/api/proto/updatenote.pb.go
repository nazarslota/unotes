// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: updatenote.proto

package proto

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type UpdateNoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	NewTitle          string                 `protobuf:"bytes,2,opt,name=new_title,json=newTitle,proto3" json:"new_title,omitempty"`
	NewContent        string                 `protobuf:"bytes,3,opt,name=new_content,json=newContent,proto3" json:"new_content,omitempty"`
	NewPriority       *string                `protobuf:"bytes,5,opt,name=new_priority,json=newPriority,proto3,oneof" json:"new_priority,omitempty"`
	NewCompletionTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=new_completion_time,json=newCompletionTime,proto3,oneof" json:"new_completion_time,omitempty"`
}

func (x *UpdateNoteRequest) Reset() {
	*x = UpdateNoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_updatenote_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateNoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateNoteRequest) ProtoMessage() {}

func (x *UpdateNoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_updatenote_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateNoteRequest.ProtoReflect.Descriptor instead.
func (*UpdateNoteRequest) Descriptor() ([]byte, []int) {
	return file_updatenote_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateNoteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateNoteRequest) GetNewTitle() string {
	if x != nil {
		return x.NewTitle
	}
	return ""
}

func (x *UpdateNoteRequest) GetNewContent() string {
	if x != nil {
		return x.NewContent
	}
	return ""
}

func (x *UpdateNoteRequest) GetNewPriority() string {
	if x != nil && x.NewPriority != nil {
		return *x.NewPriority
	}
	return ""
}

func (x *UpdateNoteRequest) GetNewCompletionTime() *timestamppb.Timestamp {
	if x != nil {
		return x.NewCompletionTime
	}
	return nil
}

type UpdateNoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateNoteResponse) Reset() {
	*x = UpdateNoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_updatenote_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateNoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateNoteResponse) ProtoMessage() {}

func (x *UpdateNoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_updatenote_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateNoteResponse.ProtoReflect.Descriptor instead.
func (*UpdateNoteResponse) Descriptor() ([]byte, []int) {
	return file_updatenote_proto_rawDescGZIP(), []int{1}
}

var File_updatenote_proto protoreflect.FileDescriptor

var file_updatenote_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaf, 0x02, 0x0a,
	0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x27, 0x0a, 0x09,
	0x6e, 0x65, 0x77, 0x5f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x00, 0x18, 0x80, 0x01, 0x52, 0x08, 0x6e, 0x65, 0x77,
	0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x2b, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72,
	0x05, 0x10, 0x00, 0x18, 0x80, 0x08, 0x52, 0x0a, 0x6e, 0x65, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x30, 0x0a, 0x0c, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xa0,
	0x01, 0x02, 0x48, 0x00, 0x52, 0x0b, 0x6e, 0x65, 0x77, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x88, 0x01, 0x01, 0x12, 0x4f, 0x0a, 0x13, 0x6e, 0x65, 0x77, 0x5f, 0x63, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x01, 0x52,
	0x11, 0x6e, 0x65, 0x77, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x42, 0x16, 0x0a, 0x14, 0x5f, 0x6e, 0x65, 0x77, 0x5f, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x14,
	0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x7a, 0x61, 0x72, 0x73, 0x6c, 0x6f, 0x74, 0x61, 0x2f, 0x75, 0x6e,
	0x6f, 0x74, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_updatenote_proto_rawDescOnce sync.Once
	file_updatenote_proto_rawDescData = file_updatenote_proto_rawDesc
)

func file_updatenote_proto_rawDescGZIP() []byte {
	file_updatenote_proto_rawDescOnce.Do(func() {
		file_updatenote_proto_rawDescData = protoimpl.X.CompressGZIP(file_updatenote_proto_rawDescData)
	})
	return file_updatenote_proto_rawDescData
}

var file_updatenote_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_updatenote_proto_goTypes = []interface{}{
	(*UpdateNoteRequest)(nil),     // 0: UpdateNoteRequest
	(*UpdateNoteResponse)(nil),    // 1: UpdateNoteResponse
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_updatenote_proto_depIdxs = []int32{
	2, // 0: UpdateNoteRequest.new_completion_time:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_updatenote_proto_init() }
func file_updatenote_proto_init() {
	if File_updatenote_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_updatenote_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateNoteRequest); i {
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
		file_updatenote_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateNoteResponse); i {
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
	file_updatenote_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_updatenote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_updatenote_proto_goTypes,
		DependencyIndexes: file_updatenote_proto_depIdxs,
		MessageInfos:      file_updatenote_proto_msgTypes,
	}.Build()
	File_updatenote_proto = out.File
	file_updatenote_proto_rawDesc = nil
	file_updatenote_proto_goTypes = nil
	file_updatenote_proto_depIdxs = nil
}
