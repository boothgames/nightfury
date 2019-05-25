package socket

import "github.com/boothgames/nightfury/pkg/db"

// Message represents message sent across ws
type Message struct {
	Action  string
	Payload db.Model
}
