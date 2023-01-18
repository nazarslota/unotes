package grpc

import (
	"context"

	pb "github.com/nazarslota/unotes/note/api/proto"
)

type noteServiceServer struct {
	services service.Services
	pb.NoteServiceServer
}

func (n noteServiceServer) CreateNote(ctx context.Context, request *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteServiceServer) UpdateNote(ctx context.Context, request *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteServiceServer) DeleteNote(ctx context.Context, request *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteServiceServer) GetNote(ctx context.Context, request *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteServiceServer) GetNotes(request *pb.GetNotesRequest, server pb.NoteService_GetNotesServer) error {
	//TODO implement me
	panic("implement me")
}
