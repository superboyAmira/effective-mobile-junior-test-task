package test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"online-song-library/internal/controller"
	"online-song-library/internal/model"
	"online-song-library/internal/router"
	external_api_test "online-song-library/test/external_api"
	mocks "online-song-library/test/mock"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(mocks.MockSongService)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	songController := controller.NewSongController(mockService, mockLogger)
	router := router.SetupRouter(songController, mockLogger)

	external_api_test.CreateMockExternalAPIServer(mockLogger)

	song := model.Song{
		Group: "Musa",
		Title: "Supermassive Black Hole",
	}
	body, _ := json.Marshal(song)
	songID := uuid.New()

	t_time, _ := time.Parse("02.01.2006", "16.07.2006")
	checkSong := model.Song{
		Id:          songID,
		Group:       "Musa",
		Title:       "Supermassive Black Hole",
		ReleaseDate: t_time,
		Text:        "Ooh baby, don't you know I suffer?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	mockService.On("UpdateSong", mock.Anything, mock.Anything, mock.Anything).Return(checkSong, nil)

	req, err := http.NewRequest(http.MethodPut, "/songs/"+songID.String(), bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedSong model.Song
	err = json.Unmarshal(w.Body.Bytes(), &updatedSong)
	assert.NoError(t, err)
	assert.Equal(t, checkSong.Group, updatedSong.Group)
	assert.Equal(t, checkSong.Title, updatedSong.Title)
	assert.Equal(t, checkSong.ReleaseDate, updatedSong.ReleaseDate)
	assert.Equal(t, checkSong.Text, updatedSong.Text)
	assert.Equal(t, checkSong.Link, updatedSong.Link)
}

func TestCreateSong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// init all mock components and controller
	mockService := new(mocks.MockSongService)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	songController := controller.NewSongController(mockService, mockLogger)
	router := router.SetupRouter(songController, mockLogger)
	external_api_test.CreateMockExternalAPIServer(mockLogger)

	// mock service settings
	tmpTime, _ := time.Parse("02.01.2006", "16.07.2006")
	mockService.On("CreateSong", mock.Anything, mock.Anything, mock.Anything).Return(uuid.New(), nil)
	mockService.On("FetchSongDetailsFromAPI", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Song{
		ReleaseDate: tmpTime,
		Text:        "Ooh baby, don't you know I suffer?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}, nil)

	// create request and req body
	song := model.Song{
		Group: "Musa",
		Title: "Supermassive Black Hole",
	}
	body, _ := json.Marshal(song)
	req, err := http.NewRequest(http.MethodPost, "/songs", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// send request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// test asserting
	assert.Equal(t, http.StatusOK, w.Code)

	response := struct {
		SongID string `json:"song_id"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	_, err = uuid.Parse(response.SongID)
	assert.NoError(t, err, "song_id is not a valid UUID")
	// assert.Equal(t,  updatedSong.Group, song.Group)
	// assert.Equal(t,updatedSong.Title, song.Title )
	// assert.Equal(t, checkTime, song.ReleaseDate)
	// assert.Equal(t, "Ooh baby, don't you know I suffer?", song.Text)
	// assert.Equal(t, "https://www.youtube.com/watch?v=Xsp3_a-PMTw", song.Link)
}

func TestDeleteSong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// init all mock components and controller
	mockService := new(mocks.MockSongService)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	songController := controller.NewSongController(mockService, mockLogger)
	router := router.SetupRouter(songController, mockLogger)
	external_api_test.CreateMockExternalAPIServer(mockLogger)

	// mock service settings
	mockService.On("DeleteSong", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// create request and req body
	uuidToDel := uuid.New()
	body, _ := json.Marshal(uuidToDel)
	req, err := http.NewRequest(http.MethodDelete, "/songs/"+uuidToDel.String(), bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// send request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// test asserting
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetLibrary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// init all mock components and controller
	mockService := new(mocks.MockSongService)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	songController := controller.NewSongController(mockService, mockLogger)
	router := router.SetupRouter(songController, mockLogger)
	external_api_test.CreateMockExternalAPIServer(mockLogger)

	// mock service settings
	tmpTime, _ := time.Parse("02.01.2006", "16.07.2006")
	tmpSong := model.Song{
		Id:          uuid.New(),
		Group:       "Musa",
		Title:       "Supermassive Black Hole",
		ReleaseDate: tmpTime,
		Text:        "Ooh baby, don't you know I suffer?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	mockService.On("GetLibrary", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]model.Song{ tmpSong }, nil)

	// create request and req body
	req, err := http.NewRequest(http.MethodGet, "/songs", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// send request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// test asserting
	assert.Equal(t, http.StatusOK, w.Code)

	checkSongs := []model.Song{}
	err = json.Unmarshal(w.Body.Bytes(), &checkSongs)
	assert.NoError(t, err)
	assert.Equal(t,  tmpSong.Id, checkSongs[0].Id)
	assert.Equal(t, tmpSong.Title, checkSongs[0].Title)
	assert.Equal(t, tmpSong.Group, checkSongs[0].Group)
	assert.Equal(t, tmpSong.ReleaseDate, checkSongs[0].ReleaseDate)
	assert.Equal(t, "Ooh baby, don't you know I suffer?", checkSongs[0].Text)
	assert.Equal(t, "https://www.youtube.com/watch?v=Xsp3_a-PMTw", checkSongs[0].Link)
}

func TestGetSongVerses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// init all mock components and controller
	mockService := new(mocks.MockSongService)
	mockLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	songController := controller.NewSongController(mockService, mockLogger)
	router := router.SetupRouter(songController, mockLogger)
	external_api_test.CreateMockExternalAPIServer(mockLogger)

	// mock service settings
	songID := uuid.New()
	mockVerses := []string{
		"First verse of the song.",
		"Second verse of the song.",
		"Third verse of the song.",
	}
	mockService.On("GetSongVerses", mock.Anything, mock.Anything, songID, 1, 5).
		Return(mockVerses, nil)

	// create request
	req, err := http.NewRequest(http.MethodGet, "/songs/"+songID.String()+"/verses?page=1&page_size=5", nil)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// send request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// test asserting
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем результат
	var returnedVerses []string
	err = json.Unmarshal(w.Body.Bytes(), &returnedVerses)
	assert.NoError(t, err)
	assert.Equal(t, mockVerses, returnedVerses)
}