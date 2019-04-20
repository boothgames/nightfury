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
	gameStarted   = "started"
	gameCompleted = "completed"
	gameFailed    = "failed"
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
	log.Infof("game '%v' of client '%v' connected", game.Name, client.Name)
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
	log.Infof("game '%v' of client '%v' disconnected", game.Name, client.Name)
}

func gameMessageReceived(session *melody.Session, data []byte) {
	client, _, err := clientFromSession(session, func(id string) (client nightfury.Client, e error) {
		return nightfury.Client{}, fmt.Errorf("client not found")
	})
	if err != nil {
		logErr(err)
		return
	}

	game, _, err := gameFromSession(session, func(name string) (nightfury.Game, error) {
		return nightfury.Game{}, fmt.Errorf("game %v not found", name)
	})
	if err != nil {
		logErr(err)
		return
	}

	message := Message{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		logErr(err)
		return
	}
	processGameMessage(message, *client, *game)
}

func processGameMessage(message Message, client nightfury.Client, game nightfury.Game) {
	switch message.Action {
	case gameStarted:
		handleGameStarted(game, client)
	case gameCompleted:
		handleGameCompleted(game, client)
	case gameFailed:
		handleGameFailed(game, client)
	default:
		err := fmt.Errorf("unknown action '%v' from game '%v' of client '%v'", message.Action, game.Name, client.Name)
		logErr(err)
	}
}

func handleGameFailed(game nightfury.Game, client nightfury.Client) {
	log.Infof("game '%v' of client '%v' has failed", game.Name, client.Name)
	if err := client.FailGame(game); err != nil {
		logErr(err)
		return
	}
	message := Message{Action: gameFailed, Payload: game}
	broadcastMessageToClient(client, message)
}

func handleGameCompleted(game nightfury.Game, client nightfury.Client) {
	log.Infof("game '%v' of client '%v' has completed playing", game.Name, client.Name)
	if err := client.CompleteGame(game); err != nil {
		logErr(err)
		return
	}
	message := Message{Action: gameCompleted, Payload: game}
	broadcastMessageToClient(client, message)

	if client.HasNext() {
		nextGame, err := client.Next()
		if err != nil {
			logErr(err)
			return
		}
		handleGameStarted(nextGame, client)
		return
	}
}

func handleGameStarted(game nightfury.Game, client nightfury.Client) {
	log.Infof("game '%v' of client '%v' has started playing", game.Name, client.Name)
	message := Message{Action: gameStarted, Payload: game}
	broadcastMessageToClient(client, message)
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
