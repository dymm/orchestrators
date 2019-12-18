package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dymm/orchestrators/grpc/cmd/processor-add/api"
)

type addServiceServer struct {
}

func (s *addServiceServer) Add(ctx context.Context, req *api.AddRequest) (*api.AddResult, error) {

	result := req.Value1 + req.Value2
	if result == 0 {
		fmt.Printf("%s - processor-add : timing out the process\n", req.Name)
		time.Sleep(4 * time.Second)
	}
	return &api.AddResult{Result: result}, nil
}
