package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dymm/orchestrators/grpc/cmd/processor-sub/api"
)

type subServiceServer struct {
}

func (s *subServiceServer) Sub(ctx context.Context, req *api.SubRequest) (*api.SubResult, error) {

	result := req.Value1 + req.Value2
	if result == 0 {
		fmt.Printf("%s - processor-sub : timing out the process\n", req.Name)
		time.Sleep(4 * time.Second)
	}
	return &api.SubResult{Result: result}, nil
}
