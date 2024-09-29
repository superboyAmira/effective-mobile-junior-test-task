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
	serv service.Service
	log  *slog.Logger
}

func NewSongController(service service.Service, logger *slog.Logger) *SongController {
	return &SongController{
		serv: service,
		log:  logger,
	}
}

// CreateSong creates a new song
// @Summary Create a new song
// @Description Creates a new song with the given details
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body model.SongDTO true "Song details"
// @Success 200 {object} model.ErrorResponse "song_id"
// @Failure 400 {object} model.ErrorResponse "Invalid input"
// @Failure 500 {object} model.ErrorResponse "Failed to create song"
// @Router /songs [post]
func (r *SongController) CreateSong(c *gin.Context) {
	var songDTO model.SongDTO
	if err := c.ShouldBindJSON(&songDTO); err != nil {
		r.log.Error("Failed to bind songDTO", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	enrichedData, err := r.serv.FetchSongDetailsFromAPI(c.Request.Context(), r.log, songDTO.Group, songDTO.Title)
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

// UpdateSong updates an existing song
// @Summary Update an existing song
// @Description Updates a song with the given ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path string true "Song ID"
// @Param song body model.Song true "Updated song details"
// @Success 200 {object} model.Song
// @Failure 400 {object} model.ErrorResponse "Invalid input or ID"
// @Failure 404 {object} model.ErrorResponse "Record not found"
// @Failure 500 {object} model.ErrorResponse "Failed to update song"
// @Router /songs/{id} [put]
func (r *SongController) UpdateSong(c *gin.Context) {
	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		r.log.Error("Failed to bind song data", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		r.log.Error("Failed to parse song ID", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song.Id = id
	updatedSong, err := r.serv.UpdateSong(c.Request.Context(), r.log, song)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		r.log.Error("Failed to update song", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, updatedSong)
}

// DeleteSong deletes a song by ID
// @Summary Delete a song
// @Description Deletes a song with the given ID
// @Tags songs
// @Param id path string true "Song ID"
// @Success 200 {object} model.ErrorResponse "Song deleted successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid song ID"
// @Failure 404 {object} model.ErrorResponse "Record not found"
// @Failure 500 {object} model.ErrorResponse "Failed to delete song"
// @Router /songs/{id} [delete]
func (r *SongController) DeleteSong(c *gin.Context) {
	songIdStr := c.Param("id")
	songId, err := uuid.Parse(songIdStr)
	if err != nil {
		r.log.Error("Invalid song ID", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	err = r.serv.DeleteSong(c.Request.Context(), r.log, songId)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}
		r.log.Error("Failed to delete song", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// GetLibrary returns a list of songs
// @Summary Get all songs in the library
// @Description Returns a list of all songs with optional filtering and pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param id query string false "Song ID"
// @Param group query string false "Group name"
// @Param title query string false "Song title"
// @Success 200 {array} model.Song
// @Failure 400 {object} model.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} model.ErrorResponse "Failed to get library"
// @Router /songs [get]
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
	if err != nil {
		r.log.Error("Invalid pagination parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		r.log.Error("Invalid pagination parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
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

// GetSongVerses returns paginated verses for a song
// @Summary Get song verses
// @Description Returns paginated verses for the specified song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path string true "Song ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of verses per page"
// @Success 200 {array} string
// @Failure 400 {object} model.ErrorResponse "Invalid song ID or pagination parameters"
// @Failure 500 {object} model.ErrorResponse "Failed to get song verses"
// @Router /songs/{id}/verses [get]
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
	if err != nil {
		r.log.Error("Invalid pagination parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		r.log.Error("Invalid pagination parameters", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
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
