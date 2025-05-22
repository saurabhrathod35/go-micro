package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WEBPORT  = "80"
	RPCPORT  = "5001"
	MONGOURL = "mongodb://mongo:27017"
	GRPCPORT = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models

	// add any configuration settings you need

}

func main() {
	// connect to mongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	log.Println("Starting service on port", WEBPORT)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEBPORT),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    ":" + WEBPORT,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(MONGOURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	return c, nil
}
