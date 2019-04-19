package socket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/log"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

// HandleGames handle socket connection related to games
func HandleGames(c *gin.Context) {
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

func gameConnected(session *melody.Session) {
	client, _, err := clientFromSession(session, func(id string) (client nightfury.Client, e error) {
		return nightfury.Client{}, fmt.Errorf("client not found")
	})
	if err != nil {
		logErr(err)
		err := session.Close()
		logErr(err)
		return
	}
	if !client.Available {
		log.Error("client not available, closing session")
		err := session.Close()
		logErr(err)
		return
	}
	game, repository, err := gameFromSession(session, func(id string) (nightfury.Game, error) {
		return nightfury.Game{Name: id}, nil
	})
	if err != nil {
		logErr(err)
		return
	}
	client.Add(*game)
	err = client.Save(repository)
	logErr(err)
}

func gameDisconnected(session *melody.Session) {
	client, _, err := clientFromSession(session, func(id string) (client nightfury.Client, e error) {
		return nightfury.Client{}, fmt.Errorf("client not found")
	})
	if err != nil {
		logErr(err)
		return
	}
	game, repository, err := gameFromSession(session, func(id string) (nightfury.Game, error) {
		return nightfury.Game{Name: id}, nil
	})
	if err != nil {
		logErr(err)
		return
	}
	client.Remove(*game)
	err = client.Save(repository)
	logErr(err)
}

func gameFromSession(session *melody.Session, notFoundFn func(name string) (nightfury.Game, error)) (*nightfury.Game, db.Repository, error) {
	if name, ok := gameName(session); ok {
		repository := db.DefaultRepository()
		game, err := nightfury.NewGameFromRepoWithName(repository, name)
		if _, ok := err.(db.EntryNotFound); ok {
			game, err := notFoundFn(name)
			if err != nil {
				return nil, nil, err
			}
			return &game, repository, nil
		} else if err != nil {
			return nil, nil, err
		}
		return &game, repository, nil
	}
	return nil, nil, fmt.Errorf("unable to parse game name from session")
}

func gameMessageReceived(session *melody.Session, msg []byte) {
	_ = gameEngine.BroadcastOthers([]byte("broadcast"), session)
}
