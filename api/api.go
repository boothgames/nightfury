package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/log"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

func init() {
	log.SetLogLevel("info")
}

// Bind binds the route to gin
func Bind(engine *gin.Engine) {
	gameEngine := melody.New()
	clientEngine := melody.New()
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	engine.GET("clients/:id/ws", func(c *gin.Context) {
		clientID := c.Param("id")
		err := clientEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"id": clientID})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	engine.GET("clients/:id/games/:name/ws", func(c *gin.Context) {
		gameName := c.Param("name")
		err := gameEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{"name": gameName})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})

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

func gameName(session *melody.Session) string {
	if session == nil {
		return ""
	}
	if name, ok := session.Keys["name"]; ok {
		return fmt.Sprintf("%v", name)
	}
	return ""
}
