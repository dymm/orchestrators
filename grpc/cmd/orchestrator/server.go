package main

import (
	"context"
	"log"
	"time"

	"github.com/dymm/orchestrators/grpc/cmd/orchestrator/api"
	add "github.com/dymm/orchestrators/grpc/cmd/processor-add/api"
	print "github.com/dymm/orchestrators/grpc/cmd/processor-print/api"
	sub "github.com/dymm/orchestrators/grpc/cmd/processor-sub/api"
	"google.golang.org/grpc"
)

type processServiceServer struct {
	infoLogger     *log.Logger
	connAdd        *grpc.ClientConn
	processorAdd   add.AddServiceClient
	connSub        *grpc.ClientConn
	processorSub   sub.SubServiceClient
	connPrint      *grpc.ClientConn
	processorPrint print.PrintServiceClient
}

func (s *processServiceServer) Process(ctx context.Context, request *api.ProcessRequest) (*api.ProcessResponse, error) {

	var result int64
	var err error
	if request.GetValue() < 50 {
		result, err = s.doProcessing1(ctx, request)
	} else {
		result, err = s.doProcessing2(ctx, request)
	}

	return &api.ProcessResponse{Result: result}, err
}

func (s *processServiceServer) doProcessing1(parentCtx context.Context, request *api.ProcessRequest) (int64, error) {
	s.infoLogger.Printf("%s - doProcessing1 start\n", request.GetName())
	defer s.infoLogger.Printf("%s - doProcessing1 end\n", request.GetName())

	//Sub
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*4)
	result, err := s.processorSub.Sub(ctx, &sub.SubRequest{Name: request.Name, Value1: request.Value, Value2: 1})
	defer cancel()
	if err != nil {
		return 0, err
	}

	//Print
	ctx2, cancel2 := context.WithTimeout(parentCtx, time.Second*4)
	defer cancel2()
	if _, err = s.processorPrint.Print(ctx2, &print.PrintRequest{Name: request.Name, Value: result.Result}); err != nil {
		return 0, err
	}

	return result.Result, nil
}

func (s *processServiceServer) doProcessing2(parentCtx context.Context, request *api.ProcessRequest) (int64, error) {
	s.infoLogger.Printf("%s - doProcessing2 start\n", request.GetName())
	defer s.infoLogger.Printf("%s - doProcessing2 end\n", request.GetName())

	//Add
	result := request.GetValue()
	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithTimeout(parentCtx, time.Second*4)
		addedResult, err := s.processorAdd.Add(ctx, &add.AddRequest{Name: request.Name, Value1: result, Value2: 1})
		defer cancel()
		if err != nil {
			return 0, err
		}
		result = addedResult.Result
	}

	//Print
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*4)
	defer cancel()
	if _, err := s.processorPrint.Print(ctx, &print.PrintRequest{Name: request.Name, Value: result}); err != nil {
		return 0, err
	}

	return result, nil
}
