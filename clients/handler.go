package clients

import (
	"encoding/json"
	"net/http"
)

func handleGet(w http.ResponseWriter, r *http.Request, model ClientModel) {
	clients, err := model.GetAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&clients)
}

func NewClientHandler(model ClientModel) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleGet(w, r, model)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	return http.HandlerFunc(fn)
}
