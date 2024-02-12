package logs

type Log struct {
	ClientId string `json:"client_id" bson:"client_id"`
	Details  string `json:"details" bson:"details"`
	Message  string `json:"message" bson:"message"`
}

func NewLog(clientId string, details string, message string) Log {
	return Log{
		ClientId: clientId,
		Details:  details,
		Message:  message,
	}
}
