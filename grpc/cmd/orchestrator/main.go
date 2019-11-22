package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dymm/orchestrators/grpc/pkg/messaging/process"
	"google.golang.org/grpc"
)

var infoLogger *log.Logger

func main() {
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		infoLogger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	connAdd, processorAdd := createProcessServiceClient("processor_add:3000")
	defer connAdd.Close()
	connSub, processorSub := createProcessServiceClient("processor_sub:3000")
	defer connSub.Close()
	connPrint, processorPrint := createProcessServiceClient("processor_print:3000")
	defer connPrint.Close()

	orchestratorServer := processServiceServer{
		connAdd:        connAdd,
		processorAdd:   processorAdd,
		connSub:        connSub,
		processorSub:   processorSub,
		connPrint:      connPrint,
		processorPrint: processorPrint,
		infoLogger :infoLogger,
	}
	process.RegisterProcessServiceServer(grpcServer, &orchestratorServer)
	if err := grpcServer.Serve(lis); err != nil {
		infoLogger.Fatalf("Failed to serve: %v", err)
	}
}

func createProcessServiceClient(addr string) (*grpc.ClientConn, process.ProcessServiceClient) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		infoLogger.Fatalf("Dial Failed to %s: %v", addr, err)
	}
	infoLogger.Println("Connected to ", addr)
	client := process.NewProcessServiceClient(conn)
	return conn, client
}
