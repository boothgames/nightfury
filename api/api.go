package api

import (
	"github.com/boothgames/nightfury/api/socket"
	"github.com/boothgames/nightfury/log"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetLogLevel("info")
}

// Bind binds the route to gin
func Bind(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := engine.Group("/v1")
	{
		v1.POST("/bulk/games", uploadGames)
		v1.POST("/bulk/hints", uploadHints)

		v1.GET("/games", listGames)
		v1.POST("/games", createGame)
		v1.GET("/games/:id", populateGame, readGame)
		v1.PUT("/games/:id", populateGame, updateGame)
		v1.DELETE("/games/:id", populateGame, deleteGame)

		v1.GET("/hints", listHints)
		v1.POST("/hints", createHint)
		v1.GET("/hints/:id", populateHint, readHint)
		v1.PUT("/hints/:id", populateHint, updateHint)
		v1.DELETE("/hints/:id", populateHint, deleteHint)

		v1.GET("/clients", listClients)
	}

	wsV1 := engine.Group("/ws/v1")
	{
		wsV1.GET("clients/:id", socket.HandleClients)
		wsV1.GET("clients/:id/games/:name", socket.HandleGames)
	}
	socket.BindSocket()
}
