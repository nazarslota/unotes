package note

import (
	"context"

	domainnote "github.com/nazarslota/unotes/note/internal/domain/note"
	"github.com/stretchr/testify/mock"
)

type mockNoteRepository struct {
	mock.Mock
}

func (m *mockNoteRepository) SaveOne(ctx context.Context, note domainnote.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *mockNoteRepository) FindOne(ctx context.Context, noteID string) (domainnote.Note, error) {
	args := m.Called(ctx, noteID)
	value, ok := args.Get(0).(domainnote.Note)
	if ok {
		return value, args.Error(1)
	}
	return domainnote.Note{}, args.Error(1)
}

func (m *mockNoteRepository) FindMany(ctx context.Context, userID string) ([]domainnote.Note, error) {
	args := m.Called(ctx, userID)
	notes, ok := args.Get(0).([]domainnote.Note)
	if ok {
		return notes, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockNoteRepository) UpdateOne(ctx context.Context, note domainnote.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *mockNoteRepository) DeleteOne(ctx context.Context, noteID string) error {
	args := m.Called(ctx, noteID)
	return args.Error(0)
}
