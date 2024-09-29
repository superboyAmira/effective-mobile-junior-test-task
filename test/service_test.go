package test

import (
	"context"
	"log/slog"
	"online-song-library/internal/model"
	"online-song-library/internal/service"
	mocks "online-song-library/test/mock"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSongService_CreateSong(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	songService := service.NewSongService(mockRepo)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	songID := uuid.New()
	mockSong := model.Song{
		Id:    songID,
		Group: "Musa",
		Title: "Supermassive Black Hole",
	}

	mockRepo.On("Create", mock.Anything, mock.Anything, mockSong).Return(songID, nil)

	result, err := songService.CreateSong(context.Background(), mockLogger, mockSong)

	assert.NoError(t, err)
	assert.Equal(t, songID, result)
}

func TestSongService_UpdateSong(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	songService := service.NewSongService(mockRepo)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	updatedSong := model.Song{
		Id:    uuid.New(),
		Group: "Musa",
		Title: "Supermassive Black Hole",
	}

	mockRepo.On("Update", mock.Anything, mock.Anything, updatedSong).Return(updatedSong, nil)

	result, err := songService.UpdateSong(context.Background(), mockLogger, updatedSong)

	assert.NoError(t, err)
	assert.Equal(t, updatedSong, result)
}

func TestSongService_DeleteSong(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	songService := service.NewSongService(mockRepo)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	songID := uuid.New()
	mockRepo.On("Delete", mock.Anything, mock.Anything, songID).Return(nil)

	err := songService.DeleteSong(context.Background(), mockLogger, songID)

	assert.NoError(t, err)
}

func TestSongService_GetLibrary(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	songService := service.NewSongService(mockRepo)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mockSongs := []model.Song{
		{
			Id:    uuid.New(),
			Group: "Musa",
			Title: "Supermassive Black Hole",
		},
	}
	filter := model.SongFilter{}
	mockRepo.On("GetAll", mock.Anything, mock.Anything, 10, 0, filter).Return(mockSongs, nil)

	result, err := songService.GetLibrary(context.Background(), mockLogger, filter, 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, mockSongs, result)
}

func TestSongService_GetSongVerses(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	songService := service.NewSongService(mockRepo)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	songID := uuid.New()
	mockText := "First verse\n\nSecond verse\n\nThird verse"
	mockRepo.On("GetVerses", mock.Anything, mock.Anything, songID).Return(mockText, nil)

	result, err := songService.GetSongVerses(context.Background(), mockLogger, songID, 1, 2)

	expectedVerses := []string{"First verse", "Second verse"}
	assert.NoError(t, err)
	assert.Equal(t, expectedVerses, result)
}
