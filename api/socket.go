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

func handleClients(c *gin.Context) {
	clientID := c.Param("id")
	err := clientEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"id": clientID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func handleGames(c *gin.Context) {
	gameName := c.Param("name")
	err := gameEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"name": gameName})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func bindSocket(engine *gin.Engine) {
	engine.GET("clients/:id/ws", handleClients)
	engine.GET("clients/:id/games/:name/ws", handleGames)

	gameEngine.HandleConnect(func(s *melody.Session) {
		_ = s.Write([]byte("welcome game"))
	})

	clientEngine.HandleConnect(func(s *melody.Session) {
		_ = s.Write([]byte("welcome client"))
	})

	gameEngine.HandleMessage(func(session *melody.Session, msg []byte) {
		_ = gameEngine.BroadcastOthers([]byte("broadcast"), session)
	})

	gameEngine.HandleMessage(func(session *melody.Session, bytes []byte) {
		message := fmt.Sprintf("%s says %s", gameName(session), string(bytes))
		_ = gameEngine.BroadcastOthers([]byte(message), session)
		log.Info(message)
	})
}
