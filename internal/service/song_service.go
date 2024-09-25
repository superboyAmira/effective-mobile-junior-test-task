package service

import (
	"context"
	"errors"
	"online-song-library/internal/model"
	"online-song-library/internal/repository"
	"log/slog"
	"strings"
	"github.com/google/uuid"
)

type SongService struct {
	repo repository.SongRepository
}

func NewSongService(r repository.SongRepository) *SongService {
	return &SongService{
		repo: r,
	}
}

func (s *SongService) CreateSong(ctx context.Context, log *slog.Logger, song model.Song) (uuid.UUID, error) {
	return s.repo.Create(ctx, log, song)
}

func (s *SongService) UpdateSong(ctx context.Context, log *slog.Logger, song model.Song) (model.Song, error) {
	return s.repo.Update(ctx, log, song)
}

func (s *SongService) DeleteSong(ctx context.Context, log *slog.Logger, songId uuid.UUID) error {
	return s.repo.Delete(ctx, log, songId)
}

func (s *SongService) GetLibrary(ctx context.Context, log *slog.Logger, filter model.SongFilter, limit, offset int) ([]model.Song, error) {
	return s.repo.GetAll(ctx, log, limit, offset, filter)
}

func (s *SongService) GetSongVerses(ctx context.Context, log *slog.Logger, songId uuid.UUID, page, pageSize int) ([]string, error) {
	text, err := s.repo.GetVerses(ctx, log, songId)
	if err != nil {
		log.Error("failed to get song text", slog.String("err", err.Error()))
		return nil, err
	}

	verses := strings.Split(text, "\n\n")

	totalVerses := len(verses)
	totalPages := (totalVerses + pageSize - 1) / pageSize

	if page < 1 || page > totalPages {
		return nil, errors.New("invalid page number")
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalVerses {
		end = totalVerses
	}

	paginatedVerses := verses[start:end]
	return paginatedVerses, nil
}
