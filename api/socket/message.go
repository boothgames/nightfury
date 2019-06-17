package socket

import "github.com/boothgames/nightfury/pkg/db"

// Message represents message sent across ws
type Message struct {
	Action  string   `json:"action"`
	Payload db.Model `json:"payload"`
}
