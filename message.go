package runtelchat

import (
	"encoding/json"
	"time"
)

//Message represents a chat message
type Message struct {
	Text      string
	From      string
	Address   string
	TimeStamp time.Time
	client    *client
}

//String is the stringer interface
func (m Message) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}
