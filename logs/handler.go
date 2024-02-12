package logs

import (
	"encoding/json"
	"net/http"

	"github.com/zulfiqarjunejo/logs-management-system/clients"
	"github.com/zulfiqarjunejo/logs-management-system/types"
)

type CreateLogRequestBody struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func handleGet(w http.ResponseWriter, r *http.Request, model LogModel) {
	logs, err := model.GetAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&logs)
}

func handlePost(w http.ResponseWriter, r *http.Request, model LogModel) {
	const clientContextKey types.ContextKey = 0
	client := r.Context().Value(clientContextKey).(clients.Client)

	var createLogRequestBody CreateLogRequestBody

	err := json.NewDecoder(r.Body).Decode(&createLogRequestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	newLog := NewLog(client.Id, createLogRequestBody.Details, createLogRequestBody.Message)
	err = model.Create(newLog)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func NewLogHandler(model LogModel) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGet(w, r, model)
		} else if r.Method == "POST" {
			handlePost(w, r, model)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	return http.HandlerFunc(fn)
}
