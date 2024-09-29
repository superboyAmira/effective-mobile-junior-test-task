package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"online-song-library/internal/model"
	"online-song-library/internal/repository"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// for mocks
type Service interface {
	CreateSong(ctx context.Context, log *slog.Logger, song model.Song) (uuid.UUID, error)
	UpdateSong(ctx context.Context, log *slog.Logger, song model.Song) (model.Song, error)
	DeleteSong(ctx context.Context, log *slog.Logger, songId uuid.UUID) error
	GetLibrary(ctx context.Context, log *slog.Logger, filter model.SongFilter, limit, offset int) ([]model.Song, error)
	GetSongVerses(ctx context.Context, log *slog.Logger, songId uuid.UUID, page, pageSize int) ([]string, error)
	FetchSongDetailsFromAPI(ctx context.Context, log *slog.Logger, group, title string) (model.Song, error)
}


type SongService struct {
	repo repository.Repository
}

func NewSongService(r repository.Repository) *SongService {
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

	log.Debug("Page Info", slog.Int("page", page), slog.Int("pageSize", pageSize))

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

	log.Debug("Page Info Slice", slog.Int("start", start), slog.Int("end", end))

	paginatedVerses := verses[start:end]

	var finalVerses []string
	for _, verse := range paginatedVerses {
		finalVerses = append(finalVerses, strings.Split(verse, "\n")...)
	}

	return finalVerses, nil
}

func (s *SongService) FetchSongDetailsFromAPI(ctx context.Context, log *slog.Logger, group, title string) (model.Song, error) {
	path := os.Getenv("PATH_EXTERNAL_API_HTTPTEST_SERVER")

	baseURL := fmt.Sprintf("%s/info", path)
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", title)

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	log.Debug("Request URL", slog.String("url", apiURL))

	resp, err := http.Get(apiURL)
	if err != nil {
		return model.Song{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Song{}, errors.New("external API status: " + strconv.Itoa(resp.StatusCode))
	}

	var songDetails model.Song

	tmp := struct {
		ReleaseDate string
		Text        string
		Link        string
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&tmp); err != nil {
		return model.Song{}, err
	}

	songDetails.ReleaseDate, err = time.Parse("02.01.2006", tmp.ReleaseDate)
	if err != nil {
		return model.Song{}, err
	}
	songDetails.Link = tmp.Link
	songDetails.Text = tmp.Text
	
	return songDetails, nil
}
