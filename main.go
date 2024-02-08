package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const PORT = ":9001"

func main() {
	mux := http.NewServeMux()

	lh := http.HandlerFunc(logsHandler)

	mux.Handle("/logs", PrintRouteInfo(lh))

	log.Printf("Starting server on port %s\n", PORT)

	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err.Error())
	}
}

type CreateLogRequestBody struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

type Log struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

var logs []Log = []Log{}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&logs)
	} else if r.Method == "POST" {
		var createLogRequestBody CreateLogRequestBody

		err := json.NewDecoder(r.Body).Decode(&createLogRequestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		log := Log{
			Message: createLogRequestBody.Message,
			Details: createLogRequestBody.Details,
		}

		logs = append(logs, log)

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
