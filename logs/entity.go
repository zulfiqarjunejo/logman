package logs

import "time"

type Log struct {
	ClientId  string    `json:"clientId" bson:"client_id"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	Details   string    `json:"details,omitempty" bson:"details,omitempty"`
	Message   string    `json:"message" bson:"message"`
}

func NewLog(clientId string, details string, message string) Log {
	return Log{
		ClientId:  clientId,
		CreatedAt: time.Now(),
		Details:   details,
		Message:   message,
	}
}
