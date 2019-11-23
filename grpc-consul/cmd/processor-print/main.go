package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dymm/orchestrators/grpc-consul/pkg/messaging/process"
	"google.golang.org/grpc"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	process.RegisterProcessServiceServer(grpcServer, &processServiceServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
