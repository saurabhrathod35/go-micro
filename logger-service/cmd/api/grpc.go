package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {

	input := req.GetLogEntry()
	// Write logs to the database
	logEntry := data.LogEntry{
		Name: input.GetName(),
		Data: input.GetData(),
	}
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}
	// return success response
	res := &logs.LogResponse{
		Result: "logged!",
	}
	return res, nil

}

func (app *Config) gRPCListen() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPCPORT))
	if err != nil {
		log.Fatal("failed to listen:", err)
		// return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})
	log.Printf("gRPC server listening on port %s", GRPCPORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
