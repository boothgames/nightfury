package socket

import (
	"fmt"
	"gitlab.com/jskswamy/nightfury/log"
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
