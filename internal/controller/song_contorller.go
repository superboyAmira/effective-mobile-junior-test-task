package controller

import (
	"log/slog"
	"net/http"
	"online-song-library/internal/model"
	"online-song-library/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SongController struct {
	serv *service.SongService
	log *slog.Logger
}

func NewSongController(service *service.SongService, logger *slog.Logger) *SongController {
	return &SongController{
		serv: service,
		log: logger,
	}
}

func (r *SongController) CreateSong(c *gin.Context) {
	var songDTO model.SongDTO
	if err := c.ShouldBindJSON(&songDTO); err != nil {
		r.log.Error("Failed to bind song DTO", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	enrichedData, err := r.serv.FetchSongDetailsFromAPI(c.Request.Context(), songDTO.Group, songDTO.Title)
	if err != nil {
		r.log.Error("Failed to fetch song details from external API", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song details"})
		return
	}

	newSong := model.Song{
		Id:          uuid.New(),
		Group:       songDTO.Group,
		Title:       songDTO.Title,
		ReleaseDate: enrichedData.ReleaseDate, 
		Text:        enrichedData.Text,    
		Link:        enrichedData.Link, 
	}

	songID, err := r.serv.CreateSong(c.Request.Context(), r.log, newSong)
	if err != nil {
		r.log.Error("Failed to create song", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"song_id": songID.String()})
}

func (r *SongController) UpdateSong(c *gin.Context) {
	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		r.log.Error("Failed to bind song data", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedSong, err := r.serv.UpdateSong(c.Request.Context(), r.log, song)
	if err != nil {
		r.log.Error("Failed to update song", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, updatedSong)
}

func (r *SongController) DeleteSong(c *gin.Context) {
	songIdStr := c.Param("id")
	songId, err := uuid.Parse(songIdStr)
	if err != nil {
		r.log.Error("Invalid song ID", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	if err := r.serv.DeleteSong(c.Request.Context(), r.log, songId); err != nil {
		r.log.Error("Failed to delete song", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

func (r *SongController) GetLibrary(c *gin.Context) {
	var filter model.SongFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		r.log.Error("Failed to bind query parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		r.log.Error("Invalid pagination parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit or offset"})
		return
	}

	songs, err := r.serv.GetLibrary(c.Request.Context(), r.log, filter, limitInt, offsetInt)
	if err != nil {
		r.log.Error("Failed to get library", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get library"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

func (r *SongController) GetSongVerses(c *gin.Context) {
	songIdStr := c.Param("id")
	songId, err := uuid.Parse(songIdStr)
	if err != nil {
		r.log.Error("Invalid song ID", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "5")

	pageInt, err := strconv.Atoi(page)
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		r.log.Error("Invalid pagination parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page or page size"})
		return
	}
	verses, err := r.serv.GetSongVerses(c.Request.Context(), r.log, songId, pageInt, pageSizeInt)
	if err != nil {
		r.log.Error("Failed to get song verses", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song verses"})
		return
	}

	c.JSON(http.StatusOK, verses)
}
