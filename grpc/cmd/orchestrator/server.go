package main

import (
	"context"
	"log"
	"time"

	"github.com/dymm/orchestrators/grpc/pkg/messaging/process"
	"google.golang.org/grpc"
)

type processServiceServer struct {
	infoLogger     *log.Logger
	connAdd        *grpc.ClientConn
	processorAdd   process.ProcessServiceClient
	connSub        *grpc.ClientConn
	processorSub   process.ProcessServiceClient
	connPrint      *grpc.ClientConn
	processorPrint process.ProcessServiceClient
}

func (s *processServiceServer) Process(ctx context.Context, request *process.ProcessRequest) (*process.ProcessResponse, error) {

	var result uint64
	var err error
	if request.GetValue() < 50 {
		result, err = s.doProcessing1(request)
	} else {
		result, err = s.doProcessing2(request)
	}

	return &process.ProcessResponse{Result: result}, err
}

func (s *processServiceServer) doProcessing1(request *process.ProcessRequest) (uint64, error) {
	s.infoLogger.Printf("%s - doProcessing1 start\n", request.GetName())
	defer s.infoLogger.Printf("%s - doProcessing1 end\n", request.GetName())

	//Sub
	s.infoLogger.Printf("%s - doProcessing1 start sub\n", request.GetName())
	newVal := &process.ProcessRequest{Name: request.GetName(), Value: request.GetValue()}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*4)
	result, err := s.processorSub.Process(ctx, newVal)
	defer cancel()
	if err != nil {
		return 0, err
	}
	newVal.Value = int64(result.GetResult())

	//Print
	s.infoLogger.Printf("%s - doProcessing1 start print\n", request.GetName())
	ctx2, cancel2 := context.WithTimeout(context.TODO(), time.Second*4)
	defer cancel2()
	if _, err = s.processorSub.Process(ctx2, newVal); err != nil {
		return 0, err
	}

	return uint64(newVal.Value), nil
}

func (s *processServiceServer) doProcessing2(request *process.ProcessRequest) (uint64, error) {
	s.infoLogger.Printf("%s - doProcessing2 start\n", request.GetName())
	defer s.infoLogger.Printf("%s - doProcessing2 end\n", request.GetName())

	//Add
	newVal := &process.ProcessRequest{Name: request.GetName(), Value: request.GetValue()}
	for i := 0; i < 2; i++ {
		s.infoLogger.Printf("%s - doProcessing2 start add %d\n", request.GetName(), i+1)
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*4)
		result, err := s.processorAdd.Process(ctx, newVal)
		defer cancel()
		if err != nil {
			return 0, err
		}
		newVal.Value = int64(result.GetResult())
	}

	//Print
	s.infoLogger.Printf("%s - doProcessing2 start print\n", request.GetName())
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*4)
	defer cancel()
	if _, err := s.processorSub.Process(ctx, newVal); err != nil {
		return 0, err
	}

	return uint64(newVal.Value), nil
}
