package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zulfiqarjunejo/logs-management-system/models"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

type CreateLogRequestBody struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

type Env struct {
	logs models.LogModel
}

func (env Env) logsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		logs, err := env.logs.GetAll()
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

		newLog := models.NewLog(createLogRequestBody.Details, createLogRequestBody.Message)
		err = env.logs.Create(newLog)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	err := godotenv.Load(".env", ".env.local")
	if err != nil {
		log.Fatalf("Error loading environment variables: %s", err.Error())
	}

	PORT := os.Getenv("PORT")
	MONGO_URI := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, mongoOptions.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatalf("MongoDB connection failed: %+v", err.Error())
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	env := &Env{
		logs: models.LogModel{
			Mongo: client,
		},
	}

	mux := http.NewServeMux()

	lh := http.HandlerFunc(env.logsHandler)
	fs := http.FileServer(http.Dir("./static"))

	mux.Handle("/api/logs", PrintRouteInfo(lh))
	mux.Handle("/", fs)

	log.Printf("Starting server on port %s\n", PORT)

	err = http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err.Error())
	}
}
