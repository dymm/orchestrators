package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dymm/orchestrators/grpc/pkg/messaging/process"
)

const valueToAdd = 7

type processServiceServer struct {
}

func (s *processServiceServer) Process(ctx context.Context, request *process.ProcessRequest) (*process.ProcessResponse, error) {

	if request.Value+valueToAdd%20 == 0 {
		fmt.Printf("%s - processorAdd : timing out the process\n", request.Name)
		time.Sleep(4 * time.Second)
	}
	return &process.ProcessResponse{Result: uint64(request.Value) + valueToAdd}, nil
}
