package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dymm/orchestrators/grpc/cmd/processor-add/api"
	"github.com/dymm/orchestrators/grpc/pkg/tracer"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

func main() {
	// initialize tracer
	actTracer, closer, err := tracer.NewTracer()
	defer closer.Close()
	if err != nil {
		log.Fatalf("Failed to create the tracer: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		log.Fatalf("faled to listen: %v", err)
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
	api.RegisterAddServiceServer(grpcServer, &addServiceServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
