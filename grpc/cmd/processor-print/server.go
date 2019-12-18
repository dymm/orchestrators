package main

import (
	"context"
	"fmt"

	"github.com/dymm/orchestrators/grpc/cmd/processor-print/api"
)

type printServiceServer struct {
}

func (s *printServiceServer) Print(ctx context.Context, req *api.PrintRequest) (*api.PrintResult, error) {

	fmt.Printf("%s - processor-print : the value %s is %d\n", req.Name, req.Name, req.Value)
	return &api.PrintResult{Result: 0}, nil
}
