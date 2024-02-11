package logs

import (
	"encoding/json"
	"net/http"
)

type LogHandler struct {
	logs LogModel
}

func NewLogHandler(model LogModel) LogHandler {
	return LogHandler{
		logs: model,
	}
}

type CreateLogRequestBody struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func (handler LogHandler) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		logs, err := handler.logs.GetAll()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&logs)
	} else if r.Method == "POST" {
		var createLogRequestBody CreateLogRequestBody

		err := json.NewDecoder(r.Body).Decode(&createLogRequestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		newLog := NewLog(createLogRequestBody.Details, createLogRequestBody.Message)
		err = handler.logs.Create(newLog)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
