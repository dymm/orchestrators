package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dymm/orchestrators/grpc-consul/pkg/messaging/process"
	resolver "github.com/nicholasjackson/grpc-consul-resolver"
	"google.golang.org/grpc"
)

var infoLogger *log.Logger

func main() {
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		infoLogger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	r := resolver.NewServiceQueryResolver("http://consul-server:8500")
	r.PollInterval = 10 * time.Second

	connAdd, processorAdd := createProcessServiceClient(r, "processor-add")
	defer connAdd.Close()
	connSub, processorSub := createProcessServiceClient(r, "processor-sub")
	defer connSub.Close()
	connPrint, processorPrint := createProcessServiceClient(r, "processor-print")
	defer connPrint.Close()

	orchestratorServer := processServiceServerImpl{
		connAdd:        connAdd,
		processorAdd:   processorAdd,
		connSub:        connSub,
		processorSub:   processorSub,
		connPrint:      connPrint,
		processorPrint: processorPrint,
		infoLogger:     infoLogger,
	}
	process.RegisterProcessServiceServer(grpcServer, &orchestratorServer)
	if err := grpcServer.Serve(lis); err != nil {
		infoLogger.Fatalf("Failed to serve: %v", err)
	}
}

func createProcessServiceClient(resolv *resolver.ConsulResolver, addr string) (*grpc.ClientConn, process.ProcessServiceClient) {
	// Create the gRPC load balancer
	lb := grpc.RoundRobin(resolv)
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBalancer(lb),
		grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		infoLogger.Fatalf("Dial Failed to %s: %v", addr, err)
	}
	infoLogger.Println("Connected to ", addr)
	client := process.NewProcessServiceClient(conn)
	return conn, client
}
