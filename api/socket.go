package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/log"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

var gameEngine = melody.New()
var clientEngine = melody.New()

const (
	socketClientID = "id"
	socketGameID   = "name"
)

func handleClients(c *gin.Context) {
	id := c.Param("id")
	err := clientEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{
		socketClientID: id,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func handleGames(c *gin.Context) {
	clientID := c.Param("id")
	gameName := c.Param("name")
	err := gameEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{
		socketGameID:   gameName,
		socketClientID: clientID,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func bindSocket() {
	clientEngine.HandleConnect(clientConnected)
	clientEngine.HandleDisconnect(clientDisconnected)

	gameEngine.HandleConnect(gameConnected)
	gameEngine.HandleDisconnect(gameDisconnected)

	gameEngine.HandleMessage(func(session *melody.Session, msg []byte) {
		_ = gameEngine.BroadcastOthers([]byte("broadcast"), session)
	})

	gameEngine.HandleMessage(func(session *melody.Session, bytes []byte) {
		if source, ok := gameName(session); ok {
			message := fmt.Sprintf("%s says %s", source, string(bytes))
			_ = gameEngine.BroadcastOthers([]byte(message), session)
			log.Info(message)
		}
	})
}
