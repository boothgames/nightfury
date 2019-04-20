package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/log"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

const (
	startClient = "start"
	resetClient = "reset"
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
	log.Infof("client %v connected", client.Name)
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
	log.Infof("client %v disconnected", client.Name)
}

func clientMessageReceived(session *melody.Session, data []byte) {
	client, _, err := clientFromSession(session, func(id string) (client nightfury.Client, e error) {
		return nightfury.Client{}, fmt.Errorf("client %v not found", client.Name)
	})
	if err != nil {
		logErr(err)
	}

	clientMessage := Message{}
	err = json.Unmarshal(data, &clientMessage)
	if err != nil {
		log.Error(err)
		return
	}

	processClientMessage(clientMessage, *client)
}

func processClientMessage(message Message, client nightfury.Client) {
	switch message.Action {
	case startClient:
		log.Infof("client '%v' has requested to start playing", client.Name)
		firstGame, err := client.Start()
		if err != nil {
			log.Error(fmt.Errorf("cannot start games of client %v. Error: %v", client.Name, err))
			return
		}
		message := Message{Action: startClient, Payload: firstGame}
		broadcastMessageToGame(client, firstGame, message)
	case resetClient:
		log.Infof("client '%v' has requested reset games", client.Name)
	default:
		err := fmt.Errorf("unknown action '%v' from client '%v'", message.Action, client.Name)
		logErr(err)
	}
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
