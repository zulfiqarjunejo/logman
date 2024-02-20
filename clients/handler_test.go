package clients

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockClientModel struct{}

func (m MockClientModel) GetAll() ([]Client, error) {
	logs := []Client{
		{
			ApiKey:      "api_key_1",
			DisplayName: "Client #1",
			Id:          "client-1",
		},
		{
			ApiKey:      "api_key_2",
			DisplayName: "Client #2",
			Id:          "client-2",
		},
	}

	return logs, nil
}

func (m MockClientModel) FindClientById(id string) (Client, error) {
	return Client{}, nil
}

func TestGetAllLogs(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/api/clients", nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	clientModel := MockClientModel{}
	clientHandler := NewClientHandler(clientModel)

	clientHandler.ServeHTTP(recorder, request)

	var clients []Client
	err = json.NewDecoder(recorder.Body).Decode(&clients)
	if err != nil {
		log.Fatalf(err.Error())
	}

	expected := 2
	got := len(clients)

	if expected != got {
		t.Errorf("expected = %v, got = %v", expected, got)
	}

	if clients[0].DisplayName != "Client #1" {
		t.Errorf("expected = %v, got = %v", "Client #1", clients[0].DisplayName)
	}
}
