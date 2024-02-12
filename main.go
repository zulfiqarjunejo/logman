package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	logHandler := logs.NewLogHandler(&logModel)
	fs := http.FileServer(http.Dir("./static"))

	// Setup MUX
	mux := http.NewServeMux()
	mux.Handle("/api/logs", PrintRouteInfo(checkApiKey(logHandler)))
	mux.Handle("/", fs)

	log.Printf("Starting server on port %s\n", PORT)

	err = http.ListenAndServe(PORT, mux)
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err.Error())
	}
}
