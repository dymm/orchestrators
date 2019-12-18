package main

import (
	"fmt"
	"log"
	"net"
	"os"

	add "github.com/dymm/orchestrators/grpc/cmd/processor-add/api"
	print "github.com/dymm/orchestrators/grpc/cmd/processor-print/api"
	sub "github.com/dymm/orchestrators/grpc/cmd/processor-sub/api"
	"github.com/dymm/orchestrators/grpc/cmd/orchestrator/api"
	"github.com/dymm/orchestrators/grpc/pkg/tracer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var infoLogger *log.Logger

func main() {
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	// initialize tracer
	actTracer, closer, err := tracer.NewTracer()
	defer closer.Close()
	if err != nil {
		infoLogger.Fatalf("Failed to create the tracer: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		infoLogger.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// add opentracing stream interceptor to chain
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(actTracer)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			// add opentracing unary interceptor to chain
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(actTracer)),
		)),
	)

	connAdd, processorAdd := createAddServiceClient("processor-add:3000", actTracer)
	defer connAdd.Close()
	connSub, processorSub := createSubServiceClient("processor-sub:3000", actTracer)
	defer connSub.Close()
	connPrint, processorPrint := createPrintServiceClient("processor-print:3000", actTracer)
	defer connPrint.Close()

	orchestratorServer := processServiceServer{
		connAdd:        connAdd,
		processorAdd:   processorAdd,
		connSub:        connSub,
		processorSub:   processorSub,
		connPrint:      connPrint,
		processorPrint: processorPrint,
		infoLogger:     infoLogger,
	}
	api.RegisterProcessServiceServer(grpcServer, &orchestratorServer)
	if err := grpcServer.Serve(lis); err != nil {
		infoLogger.Fatalf("Failed to serve: %v", err)
	}
}

func createAddServiceClient(addr string, aTracer opentracing.Tracer) (*grpc.ClientConn, add.AddServiceClient) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(aTracer, otgrpc.LogPayloads())),
	)
	if err != nil {
		infoLogger.Fatalf("Dial Failed to %s: %v", addr, err)
	}
	infoLogger.Println("Connected to ", addr)
	client := add.NewAddServiceClient(conn)
	return conn, client
}
func createSubServiceClient(addr string, aTracer opentracing.Tracer) (*grpc.ClientConn, sub.SubServiceClient) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(aTracer, otgrpc.LogPayloads())),
	)
	if err != nil {
		infoLogger.Fatalf("Dial Failed to %s: %v", addr, err)
	}
	infoLogger.Println("Connected to ", addr)
	client := sub.NewSubServiceClient(conn)
	return conn, client
}
func createPrintServiceClient(addr string, aTracer opentracing.Tracer) (*grpc.ClientConn, print.PrintServiceClient) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(aTracer, otgrpc.LogPayloads())),
	)
	if err != nil {
		infoLogger.Fatalf("Dial Failed to %s: %v", addr, err)
	}
	infoLogger.Println("Connected to ", addr)
	client := print.NewPrintServiceClient(conn)
	return conn, client
}
