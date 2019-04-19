package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/jskswamy/nightfury/pkg/db"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

func listClients(c *gin.Context) {
	repository := db.DefaultRepository()
	clients, err := nightfury.NewClientsFromRepo(repository)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clients)
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
