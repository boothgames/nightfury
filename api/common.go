package api

import (
	"fmt"
	"gopkg.in/olahol/melody.v1"
)

func gameName(session *melody.Session) string {
	if session == nil {
		return ""
	}
	if name, ok := session.Keys["name"]; ok {
		return fmt.Sprintf("%v", name)
	}
	return ""
}
