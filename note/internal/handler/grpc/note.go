package grpc

import (
	"context"
	"errors"

	pb "github.com/nazarslota/unotes/note/api/proto"
	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/nazarslota/unotes/note/internal/service"
	servicenote "github.com/nazarslota/unotes/note/internal/service/note"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type noteServiceServer struct {
	services service.Services
	logger   Logger
	pb.NoteServiceServer
}

func (s noteServiceServer) CreateNote(ctx context.Context, in *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &servicenote.CreateNoteRequest{
		Title:   in.Title,
		Content: in.Content,
		UserID:  in.UserId,
	}

	response, err := s.services.NoteService.CreateNoteRequestHandler.Handle(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.CreateNoteResponse{Id: response.ID}, nil
}

func (s noteServiceServer) UpdateNote(_ context.Context, _ *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented")
}

func (s noteServiceServer) DeleteNote(ctx context.Context, in *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &servicenote.DeleteNoteRequest{
		ID: in.Id,
	}

	_, err := s.services.NoteService.DeleteNoteRequestHandler.Handle(ctx, request)
	if errors.Is(err, domainnote.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.DeleteNoteResponse{}, nil
}

func (s noteServiceServer) GetNote(ctx context.Context, in *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := &servicenote.GetNoteRequest{
		ID: in.Id,
	}

	response, err := s.services.NoteService.GetNoteRequestHandler.Handle(ctx, request)
	if errors.Is(err, domainnote.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.GetNoteResponse{
		Title:   response.Title,
		Content: response.Content,
		UserId:  response.UserID,
	}, nil
}

func (s noteServiceServer) GetNotes(in *pb.GetNotesRequest, server pb.NoteService_GetNotesServer) error {
	if err := in.Validate(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	request := &servicenote.GetNotesRequest{
		UserID: in.UserId,
	}

	response, err := s.services.NoteService.GetNotesRequestHandler.Handle(server.Context(), request)
	if errors.Is(err, domainnote.ErrNotFound) {
		return status.Error(codes.NotFound, "not found")
	} else if err != nil {
		return status.Error(codes.Internal, "internal")
	}

	for _, note := range response.Notes {
		err := server.Send(&pb.GetNotesResponse{
			Id:      note.ID,
			Title:   note.Title,
			Content: note.Content,
		})

		if err != nil {
			return status.Error(codes.Internal, "internal")
		}
	}
	return nil
}
