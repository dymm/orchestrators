package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dymm/orchestrators/grpc/pkg/messaging/process"
	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	process.RegisterProcessServiceServer(grpcServer, &processServiceServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
