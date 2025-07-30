package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
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

	// register RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic("error while register rpc", err)
	}
	go app.rpcListen()

	go app.gRPCListen()
	// start web server
	log.Println("Starting service on port", WEBPORT)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEBPORT),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("Error while listen", err)
	}
}

func (app *Config) rpcListen() {
	log.Println("Starting RPC server on port", RPCPORT)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", RPCPORT))
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()
	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(rpcConn)
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
