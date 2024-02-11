package logs

type Log struct {
	Details string `json:"details"`
	Message string `json:"message"`
}

func NewLog(details string, message string) Log {
	return Log{
		Details: details,
		Message: message,
	}
}
