package grpc

import (
	"context"
	"errors"

	pb "github.com/nazarslota/unotes/note/api/proto"
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
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context done")
	default:
	}

	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := servicenote.CreateNoteRequest{
		ID:      in.Id,
		Title:   in.Title,
		Content: in.Content,
		UserID:  in.UserId,
	}

	_, err := s.services.NoteService.CreateNoteRequestHandler.Handle(ctx, request)
	if errors.Is(err, servicenote.ErrCreateNoteAlreadyExist) {
		return nil, status.Error(codes.AlreadyExists, "note already exist")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.CreateNoteResponse{}, nil
}

func (s noteServiceServer) GetNote(ctx context.Context, in *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context done")
	default:
	}

	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := servicenote.GetNoteRequest{
		ID: in.Id,
	}

	response, err := s.services.NoteService.GetNoteRequestHandler.Handle(ctx, request)
	if errors.Is(err, servicenote.ErrGetNoteNotFound) {
		return nil, status.Error(codes.NotFound, "note not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.GetNoteResponse{Title: response.Title, Content: response.Content, UserId: response.UserID}, nil
}

func (s noteServiceServer) GetNotes(in *pb.GetNotesRequest, server pb.NoteService_GetNotesServer) error {
	select {
	case <-server.Context().Done():
		return status.Error(codes.Canceled, "context done")
	default:
	}

	if err := in.Validate(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	request := servicenote.GetNotesAsyncRequest{
		UserID: in.UserId,
	}

	response, errs := s.services.NoteService.GetNotesAsyncRequestHandler.Handle(server.Context(), request)
	for {
		select {
		case <-server.Context().Done():
			return status.Error(codes.Canceled, "context done")
		default:
		}

		select {
		case note, ok := <-response.Notes:
			if !ok {
				continue
			}

			if err := server.Send(&pb.GetNotesResponse{
				Id:      note.ID,
				Title:   note.Title,
				Content: note.Content,
				UserId:  note.UserID,
			}); err != nil {
				return status.Errorf(codes.Unknown, "failed to send response")
			}
		case err, ok := <-errs:
			if !ok {
				return nil
			}

			if errors.Is(err, servicenote.ErrGetNotesAsyncNotFound) {
				return status.Errorf(codes.NotFound, "notes not found")
			} else if err != nil {
				return status.Error(codes.Internal, "internal")
			}
		}
	}
}

func (s noteServiceServer) UpdateNote(ctx context.Context, in *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context done")
	default:
	}

	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := servicenote.UpdateNoteRequest{
		ID:         in.Id,
		NewTitle:   in.NewTitle,
		NewContent: in.NewContent,
	}

	_, err := s.services.NoteService.UpdateNoteRequestHandler.Handle(ctx, request)
	if errors.Is(err, servicenote.ErrUpdateNoteNotFound) {
		return nil, status.Error(codes.NotFound, "note not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.UpdateNoteResponse{}, nil
}

func (s noteServiceServer) DeleteNote(ctx context.Context, in *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "context done")
	default:
	}

	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := servicenote.DeleteNoteRequest{
		ID: in.Id,
	}

	_, err := s.services.NoteService.DeleteNoteRequestHandler.Handle(ctx, request)
	if errors.Is(err, servicenote.ErrDeleteNoteNotFound) {
		return nil, status.Error(codes.NotFound, "note not found")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "internal")
	}
	return &pb.DeleteNoteResponse{}, nil
}
