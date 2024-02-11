package logs

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockLogModel struct{}

func (m MockLogModel) GetAll() ([]Log, error) {
	logs := []Log{
		{
			Message: "first message",
			Details: "details of first message",
		},
		{
			Message: "second message",
			Details: "details of second message",
		},
	}

	return logs, nil
}

func (m MockLogModel) Create(log Log) error {
	return nil
}

func TestGetAllLogs(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/api/logs", nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	logModel := MockLogModel{}
	logHandler := NewLogHandler(logModel)

	logHandler.ServeHTTP(recorder, request)

	var logs []Log
	err = json.NewDecoder(recorder.Body).Decode(&logs)
	if err != nil {
		log.Fatalf(err.Error())
	}

	expected := 2
	got := len(logs)

	if expected != got {
		t.Errorf("expected = %v, got = %v", expected, got)
	}

	if logs[0].Message != "first message" {
		t.Errorf("expected = %v, got = %v", "first message", logs[0].Message)
	}
}
