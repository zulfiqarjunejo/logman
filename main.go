package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zulfiqarjunejo/logs-management-system/clients"
	"github.com/zulfiqarjunejo/logs-management-system/logs"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

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

	// Initialize models.
	logModel := logs.NewMongoLogModel(client)
	clientModel := clients.NewMongoClientModel(client)

	// Initialize middlewares, if needed.
	checkApiKey := CreateCheckApiKey(&clientModel)

	// Initialize handlers.
	clientHandler := clients.NewClientHandler(&clientModel)
	logHandler := logs.NewLogHandler(&logModel)
	staticFileSystemHandler := http.FileServer(http.Dir("./static"))
	swaggerHandler := http.FileServer(http.Dir("./swagger"))

	// Setup MUX
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/api/clients", PrintRouteInfo(checkApiKey(clientHandler)))
	mux.Handle("/api/logs", PrintRouteInfo(checkApiKey(logHandler)))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", swaggerHandler))

	mux.HandleFunc("/candies/", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			id := strings.TrimPrefix(r.URL.Path, "/candies/")
			if id == "" {
				// handle list operation e.g., list all candies
				candies := []string{"all", "of", "my", "candies"}
				json.NewEncoder(w).Encode(candies)
			} else {
				candy := "only one candy"
				json.NewEncoder(w).Encode(candy)
			}
		} else if method == "POST" {
			// handle post operation e.g., create a candy.
		} else {
			// handle method not allowed e.g., log it somewhere.
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/", staticFileSystemHandler)

	log.Printf("Starting server on port %s\n", PORT)

	err = http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err.Error())
	}
}
