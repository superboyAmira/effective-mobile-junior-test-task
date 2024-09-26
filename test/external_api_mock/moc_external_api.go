package test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func CreateMockExternalAPIServer(log *slog.Logger) *httptest.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")

		log.Debug("REQUEST TO MOCK", slog.String("group", group), slog.String("song", song))

		if group == "Muse" && song == "Supermassive Black Hole" {
			response := SongDetail{
				ReleaseDate: "16.07.2006",
				Text:        "Ooh baby, don't you know I suffer?",
				Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
			}
			c.JSON(http.StatusOK, response)
		} else if group == "Enigma" && song == "Sadeness" {
			c.JSON(200, SongDetail{
				ReleaseDate: "01.10.1990",
				Text: "Procedamus in pace\n" +
					"In nomine Christi, amen\n" +
					"Cum angelis et pueris\n" +
					"Fideles inveniamur\n\n" +
					"Attollite portas, principes, vestras\n" +
					"Et elevamini, portae aeternales\n" +
					"Et introibit rex gloriae\n" +
					"Quis est iste rex gloriae?\n\n" +
					"Sade, dis-moi\n" +
					"Sade, donne-moi\n\n" +
					"Procedamus in pace\n" +
					"In nomine Christi, amen\n\n" +
					"Sade, dis-moi\n" +
					"Qu'est-ce que tu vas chercher?\n" +
					"Le bien par le mal?\n" +
					"La vertu par le vice?\n\n" +
					"Sade, dis-moi\n" +
					"Pourquoi l'évangile du mal?\n" +
					"Quelle est ta religion? Où sont tes fidèles?\n" +
					"Si tu es contre Dieu, tu es contre l'homme\n" +
					"Sade, es-tu diabolique ou divin?\n\n" +
					"Sade, dis-moi (Hosanna)\n" +
					"Sade, donne-moi (Hosanna)\n" +
					"Sade, dis-moi (Hosanna)\n" +
					"Sade, donne-moi (Hosanna)\n\n" +
					"In nomine Christi, amen",
				Link: "https://www.youtube.com/watch?v=4F9DxYhqmKw&ab_channel=EnigmaVEVO",
			})
		} else if group == "Axel F" && song == "Crazy Frog" {
			c.JSON(200, SongDetail{
				ReleaseDate: "17.05.2005",
				Text: "Ring ding ding daa baa\n" +
					"Baa aramba baa bom baa barooumba\n" +
					"Wh-wha-what's going on-on?\n" +
					"Ding, ding\n" +
					"This is the Crazy Frog\n\n" +
					"Ding, ding\n" +
					"Bem bem!\n\n" +
					"Ring ding ding ding ding ding\n" +
					"Ring ding ding ding bem bem bem\n" +
					"Ring ding ding ding ding ding\n" +
					"Ring ding ding ding baa baa\n" +
					"Ring ding ding ding ding ding\n" +
					"Ring ding ding ding bem bem bem\n" +
					"Ring ding ding ding ding ding\n" +
					"This is the Crazy Frog\n" +
					"Breakdown!\n\n" +
					"Ding ding\n" +
					"Br-br-break it, br-break it\n" +
					"Dum dum dumda dum dum dum\n" +
					"Dum dum dumda dum dum dum\n" +
					"Dum dum dumda dum dum dum\n" +
					"Bem, bem!\n" +
					"Dum dum dumda dum dum dum\n" +
					"Dum dum dumda dum dum dum\n" +
					"Dum dum dumda dum dum dum\n" +
					"This is the Crazy Frog",
				Link: "https://www.youtube.com/watch?v=k85mRPqvMbE&ab_channel=CrazyFrog",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Song not found"})
		}
	})

	return httptest.NewServer(router)
}
