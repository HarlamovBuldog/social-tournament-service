package main

import (
	"context"
	"flag"
	v1 "github.com/HarlamovBuldog/social-tournament-service/internal/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func main() {
	// get configuration.
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewTournamentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createUserReq := &v1.CreateUserRequest{Name: "Stas"}
	createUserResp, err := c.CreateUser(ctx, createUserReq)
	if err != nil {
		log.Fatalf("create user req failed: %v", err)
	}

	userID := createUserResp.Id
	log.Printf("user created with id: %s\n", userID)
}
