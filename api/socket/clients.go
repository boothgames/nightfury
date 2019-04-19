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

// HandleClients handle socket connection related to clients
func HandleClients(c *gin.Context) {
	id := c.Param("id")
	err := clientEngine.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{
		socketClientID: id,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func clientConnected(session *melody.Session) {
	client, repository, err := clientFromSession(session, func(id string) (nightfury.Client, error) {
		return nightfury.NewClient(id, true), nil
	})
	if err != nil {
		logErr(err)
		return
	}
	connectedClient := client.Connected()
	err = connectedClient.Save(repository)
	logErr(err)
}

func clientDisconnected(session *melody.Session) {
	client, repository, err := clientFromSession(session, func(id string) (nightfury.Client, error) {
		return nightfury.NewClient(id, false), nil
	})
	if err != nil {
		logErr(err)
		return
	}
	connectedClient := client.Disconnected()
	err = connectedClient.Save(repository)
	logErr(err)
}

func clientFromSession(session *melody.Session, notFoundFn func(id string) (nightfury.Client, error)) (*nightfury.Client, db.Repository, error) {
	if id, ok := clientID(session); ok {
		repository := db.DefaultRepository()
		client, err := nightfury.NewClientFromRepoWithName(repository, id)
		if _, ok := err.(db.EntryNotFound); ok {
			client, err := notFoundFn(id)
			if err != nil {
				return nil, nil, err
			}
			return &client, repository, nil
		} else if err != nil {
			return nil, nil, err
		}
		return &client, repository, nil
	}
	return nil, nil, fmt.Errorf("unable to parse client id from session")
}

func clientMessageReceived(session *melody.Session, msg []byte) {
	if source, ok := gameName(session); ok {
		message := fmt.Sprintf("%s says %s", source, string(msg))
		_ = gameEngine.BroadcastOthers([]byte(message), session)
		log.Info(message)
	}
}
