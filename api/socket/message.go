package socket

import "gitlab.com/jskswamy/nightfury/pkg/db"

// Message represents message sent across ws
type Message struct {
	Action  string
	Payload db.Model
}
