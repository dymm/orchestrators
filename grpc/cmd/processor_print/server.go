package main

import (
	"context"
	"fmt"

	"github.com/dymm/orchestrators/grpc/pkg/messaging/process"
)

const valueToAdd = 7

type processServiceServer struct {
}

func (s *processServiceServer) Process(ctx context.Context, request *process.ProcessRequest) (*process.ProcessResponse, error) {

	fmt.Printf("%s - value is %d\n", request.Name, request.Value)

	return &process.ProcessResponse{}, nil
}
