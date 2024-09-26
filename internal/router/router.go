package router

import (
	"log/slog"
	"online-song-library/internal/controller"

	"github.com/gin-gonic/gin"
)


func SetupRouter(songController *controller.SongController, log *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		log.Debug("INCOMING REQUEST", slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path) , slog.String("IP", c.ClientIP()))
		c.Next()
	})

	router.POST("/songs", songController.CreateSong)
	router.PUT("/songs/:id", songController.UpdateSong)
	router.DELETE("/songs/:id", songController.DeleteSong)
	router.GET("/songs", songController.GetLibrary)
	router.GET("/songs/:id/verses", songController.GetSongVerses)

	return router
}