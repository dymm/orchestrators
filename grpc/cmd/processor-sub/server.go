package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dymm/orchestrators/grpc/pkg/messaging/process"
)

const valueToAdd = -7

type processServiceServer struct {
}

func (s *processServiceServer) Process(ctx context.Context, request *process.ProcessRequest) (*process.ProcessResponse, error) {

	if request.Value-valueToAdd%33 == 0 {
		fmt.Printf("%s - processor-sub : timing out the process\n", request.Name)
		time.Sleep(4 * time.Second)
	}

	result := request.Value + valueToAdd
	if result < 0 {
		result = 0
	}
	return &process.ProcessResponse{Result: uint64(result)}, nil
}
