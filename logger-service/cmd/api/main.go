package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PORT      = "80"
	RPC_PORT  = "5001"
	MONGO_URL = "mongodb://mongo:27017"
	GRPC_PORT = "50001"
)

var client *mongo.Client

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: *data.New(client),
	}

	app.Serve()
}

func (cfg *Config) Serve() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: cfg.routes(),
	}

	log.Println("logger listening on port: ", PORT)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	clientOptions.SetAuth(options.Credential{
		Username: "mongo",
		Password: "password",
	})

	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Connected to Mongo successfully")
	return mongoClient, nil
}
