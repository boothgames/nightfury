package socket

import (
	"encoding/json"
	"fmt"
	"gitlab.com/jskswamy/nightfury/log"
	"gitlab.com/jskswamy/nightfury/pkg/nightfury"
	"gopkg.in/olahol/melody.v1"
)

var gameEngine = melody.New()
var clientEngine = melody.New()

const (
	socketClientID = "id"
	socketGameID   = "name"
)

// BindSocket binds the necessary sockets related to games and clients
func BindSocket() {
	clientEngine.HandleConnect(clientConnected)
	clientEngine.HandleDisconnect(clientDisconnected)
	clientEngine.HandleMessage(clientMessageReceived)

	gameEngine.HandleConnect(gameConnected)
	gameEngine.HandleDisconnect(gameDisconnected)
	gameEngine.HandleMessage(gameMessageReceived)
}

func gameName(session *melody.Session) (string, bool) {
	if session == nil {
		return "", false
	}
	if name, ok := session.Keys[socketGameID]; ok {
		return fmt.Sprintf("%v", name), true
	}
	return "", false
}

func clientID(session *melody.Session) (string, bool) {
	if session == nil {
		return "", false
	}
	if name, ok := session.Keys[socketClientID]; ok {
		return fmt.Sprintf("%v", name), true
	}
	return "", false
}

func logErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

func broadcastMessageToClient(message Message, client nightfury.Client) {
	broadcastMessage(clientEngine, message, func(session *melody.Session) bool {
		if name, ok := clientID(session); ok {
			return name == client.Name
		}
		return false
	})
}

func broadcastMessageToGame(message Message, game nightfury.Game) {
	broadcastMessage(gameEngine, message, func(session *melody.Session) bool {
		if name, ok := gameName(session); ok {
			return name == game.Name
		}
		return false
	})
}

func broadcastMessage(engine *melody.Melody, message Message, predicateFn func(session *melody.Session) bool) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
		return
	}

	err = engine.BroadcastFilter(data, predicateFn)
	if err != nil {
		log.Error(err)
		return
	}
}
