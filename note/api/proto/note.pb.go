// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: note.proto

package proto

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_note_proto protoreflect.FileDescriptor

var file_note_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x67, 0x65,
	0x74, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x67, 0x65, 0x74,
	0x6e, 0x6f, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32,
	0x83, 0x03, 0x0a, 0x0b, 0x4e, 0x6f, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4b, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x22, 0x09,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x6f, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x44, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x12, 0x0f, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x10, 0x12, 0x0e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x6f, 0x74, 0x65, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x12, 0x45, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x10,
	0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x11, 0x2e, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x6e, 0x6f, 0x74, 0x65, 0x73, 0x30, 0x01, 0x12, 0x4b, 0x0a, 0x0a, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x1a, 0x09, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e,
	0x6f, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x4d, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4e, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4e, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x4e, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x10, 0x2a, 0x0e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e, 0x6f, 0x74, 0x65,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x7a, 0x61, 0x72, 0x73, 0x6c, 0x6f, 0x74, 0x61, 0x2f, 0x75,
	0x6e, 0x6f, 0x74, 0x65, 0x73, 0x2f, 0x6e, 0x6f, 0x74, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_note_proto_goTypes = []interface{}{
	(*CreateNoteRequest)(nil),  // 0: CreateNoteRequest
	(*GetNoteRequest)(nil),     // 1: GetNoteRequest
	(*GetNotesRequest)(nil),    // 2: GetNotesRequest
	(*UpdateNoteRequest)(nil),  // 3: UpdateNoteRequest
	(*DeleteNoteRequest)(nil),  // 4: DeleteNoteRequest
	(*CreateNoteResponse)(nil), // 5: CreateNoteResponse
	(*GetNoteResponse)(nil),    // 6: GetNoteResponse
	(*GetNotesResponse)(nil),   // 7: GetNotesResponse
	(*UpdateNoteResponse)(nil), // 8: UpdateNoteResponse
	(*DeleteNoteResponse)(nil), // 9: DeleteNoteResponse
}
var file_note_proto_depIdxs = []int32{
	0, // 0: NoteService.CreateNote:input_type -> CreateNoteRequest
	1, // 1: NoteService.GetNote:input_type -> GetNoteRequest
	2, // 2: NoteService.GetNotes:input_type -> GetNotesRequest
	3, // 3: NoteService.UpdateNote:input_type -> UpdateNoteRequest
	4, // 4: NoteService.DeleteNote:input_type -> DeleteNoteRequest
	5, // 5: NoteService.CreateNote:output_type -> CreateNoteResponse
	6, // 6: NoteService.GetNote:output_type -> GetNoteResponse
	7, // 7: NoteService.GetNotes:output_type -> GetNotesResponse
	8, // 8: NoteService.UpdateNote:output_type -> UpdateNoteResponse
	9, // 9: NoteService.DeleteNote:output_type -> DeleteNoteResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_note_proto_init() }
func file_note_proto_init() {
	if File_note_proto != nil {
		return
	}
	file_createnote_proto_init()
	file_getnote_proto_init()
	file_getnotes_proto_init()
	file_updatenote_proto_init()
	file_deletenote_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_note_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_note_proto_goTypes,
		DependencyIndexes: file_note_proto_depIdxs,
	}.Build()
	File_note_proto = out.File
	file_note_proto_rawDesc = nil
	file_note_proto_goTypes = nil
	file_note_proto_depIdxs = nil
}
