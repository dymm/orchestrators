package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"

	"github.com/dymm/orchestrators/grpc/cmd/orchestrator/api"
	"github.com/dymm/orchestrators/grpc/pkg/tracer"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

const orchestrator = "orchestrator:3000"
const maxValue = 999999

var infoLogger *log.Logger

func main() {
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	// initialize tracer
	actTracer, closer, err := tracer.NewTracer()
	defer closer.Close()
	if err != nil {
		infoLogger.Fatalf("Failed to create the tracer: %v", err)
	}

	conn, err := grpc.Dial(orchestrator, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(actTracer, otgrpc.LogPayloads())),
	)

	if err != nil {
		infoLogger.Fatalf("Dial Failed: %v", err)
	}
	infoLogger.Println("Connected to ", orchestrator)
	defer conn.Close()
	processClient := api.NewProcessServiceClient(conn)

	for i := 10; i > 0; i-- {
		infoLogger.Println("Starting in ", i)
		time.Sleep(1 * time.Second)
	}

	counter := 0
	var inFligth int64

	for {
		counter = counter + 1
		newValue := api.ProcessRequest{
			Name:  fmt.Sprintf("Value %d", counter),
			Value: int64(counter % maxValue),
		}
		go func(val *api.ProcessRequest) {
			atomic.AddInt64(&inFligth, 1)
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
			defer cancel()
			_, err := processClient.Process(ctx, val)
			elapsed := time.Since(start)
			if err == nil {
				infoLogger.Printf("%s - End of processing in %s\n", val.GetName(), elapsed)
			} else {
				infoLogger.Printf("%s - Error while processing in %s\n%v\n", val.GetName(), elapsed, err)
			}

			atomic.AddInt64(&inFligth, -1)
		}(&newValue)

		if counter%maxValue == 0 {
			for i := 10; i > 0; i-- {
				infoLogger.Printf("%d call pending, emiting agin in %d sec", inFligth, i)
				time.Sleep(1 * time.Second)
			}
		} else if inFligth > 200 {
			infoLogger.Printf("%d call pending, waiting a little bit", inFligth)
			for i := 5; i > 0; i-- {
				infoLogger.Printf("%d call still pending, emiting agin in %d sec", inFligth, i)
				time.Sleep(1 * time.Second)
			}
		}
	}
}
