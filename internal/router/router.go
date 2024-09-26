package router

import (
	"online-song-library/internal/controller"

	"github.com/gin-gonic/gin"
)


func SetupRouter(songController *controller.SongController) *gin.Engine {
	router := gin.Default()

	router.POST("/songs", songController.CreateSong)
	router.PUT("/songs/:id", songController.UpdateSong)
	router.DELETE("/songs/:id", songController.DeleteSong)
	router.GET("/songs", songController.GetLibrary)
	router.GET("/songs/:id/verses", songController.GetSongVerses)

	return router
}