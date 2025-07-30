package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (s *RPCServer) LogInfo(payload RPCPayload, reply *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry:", err)
		return err
	}
	log.Println("Log entry inserted:", payload.Name, payload.Data)
	*reply = "Logged successfully via RPC: " + payload.Name + " - " + payload.Data
	return nil
}
