package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/api/socket"
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

		v1.GET("/security-incidents", listSecurityIncidents)
		v1.POST("/security-incidents", createSecurityIncident)
		v1.GET("/security-incidents/:id", populateSecurityIncident, readSecurityIncident)
		v1.PUT("/security-incidents/:id", populateSecurityIncident, updateSecurityIncident)
		v1.DELETE("/security-incidents/:id", populateSecurityIncident, deleteSecurityIncident)

		v1.GET("/clients", listClients)
	}

	wsV1 := engine.Group("/ws/v1")
	{
		wsV1.GET("clients/:id", socket.HandleClients)
		wsV1.GET("clients/:id/games/:name", socket.HandleGames)
	}
	socket.BindSocket()
}
