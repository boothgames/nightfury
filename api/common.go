package api

import (
	"fmt"
	"gitlab.com/jskswamy/nightfury/log"
	"gopkg.in/olahol/melody.v1"
)

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
