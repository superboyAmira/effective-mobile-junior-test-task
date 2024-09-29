package mock
import (
	"context"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"online-song-library/internal/model"
	"github.com/google/uuid"
)

// MockSongService is a mock implementation of the SongService for testing purposes.
type MockSongService struct {
	mock.Mock
}

func (m *MockSongService) CreateSong(ctx context.Context, log *slog.Logger, song model.Song) (uuid.UUID, error) {
	args := m.Called(ctx, log, song)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockSongService) UpdateSong(ctx context.Context, log *slog.Logger, song model.Song) (model.Song, error) {
	args := m.Called(ctx, log, song)
	return args.Get(0).(model.Song), args.Error(1)
}

func (m *MockSongService) DeleteSong(ctx context.Context, log *slog.Logger, songId uuid.UUID) error {
	args := m.Called(ctx, log, songId)
	return args.Error(0)
}

func (m *MockSongService) GetLibrary(ctx context.Context, log *slog.Logger, filter model.SongFilter, limit, offset int) ([]model.Song, error) {
	args := m.Called(ctx, log, filter, limit, offset)
	return args.Get(0).([]model.Song), args.Error(1)
}

func (m *MockSongService) GetSongVerses(ctx context.Context, log *slog.Logger, songId uuid.UUID, page, pageSize int) ([]string, error) {
	args := m.Called(ctx, log, songId, page, pageSize)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockSongService) FetchSongDetailsFromAPI(ctx context.Context, log *slog.Logger, group, title string) (model.Song, error) {
	args := m.Called(ctx, log, group, title)
	return args.Get(0).(model.Song), args.Error(1)
}
