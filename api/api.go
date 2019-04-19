package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/log"
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
		v1.GET("/games", listGames)
		v1.POST("/games", createGame)
		v1.GET("/games/:id", populateGame, readGame)
		v1.PUT("/games/:id", populateGame, updateGame)
		v1.DELETE("/games/:id", populateGame, deleteGame)
		v1.GET("/clients", listClients)
	}
	bindSocket(engine)
}
