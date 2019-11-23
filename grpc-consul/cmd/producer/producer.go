package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	resolver "github.com/nicholasjackson/grpc-consul-resolver"
	"google.golang.org/grpc"

	"github.com/dymm/orchestrators/grpc-consul/pkg/messaging/process"
)

const maxValue = 200

var infoLogger *log.Logger

func main() {
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	for i := 5; i > 0; i-- {
		infoLogger.Println("Configuration in ", i)
		time.Sleep(1 * time.Second)
	}

	r := resolver.NewServiceQueryResolver("http://consul-server:8500")
	// use the default poll interval of 60 seconds
	// the poll interval can be changed by setting the resolvers PollInterval field
	r.PollInterval = 10 * time.Second

	// Create the gRPC load balancer
	lb := grpc.RoundRobin(r)

	// create a new gRPC client connection
	conn, err := grpc.Dial(
		"orchestrator",
		grpc.WithInsecure(),
		grpc.WithBalancer(lb),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		infoLogger.Fatalf("Dial Failed: %v\n", err)
	}
	infoLogger.Println("Connected to the orchestrator")
	defer conn.Close()
	processClient := process.NewProcessServiceClient(conn)

	for i := 5; i > 0; i-- {
		infoLogger.Println("Starting in ", i)
		time.Sleep(1 * time.Second)
	}

	counter := 0
	var inFligth int64

	for {
		counter = counter + 1
		newValue := process.ProcessRequest{
			Name:  fmt.Sprintf("Value %d", counter),
			Value: int64(counter % maxValue),
		}
		go func(val *process.ProcessRequest) {
			atomic.AddInt64(&inFligth, 1)
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
			defer cancel()
			_, err := processClient.Process(ctx, val)
			elapsed := time.Since(start)
			if err == nil {
				infoLogger.Printf("%s - End of processing in %s\n", val.GetName(), elapsed)
			} else {
				infoLogger.Printf("%s - Error while processing in %s\n%v\n", val.GetName(), err, elapsed)
			}

			atomic.AddInt64(&inFligth, -1)
		}(&newValue)

		if counter%maxValue == 0 {
			infoLogger.Printf("%d call pending\n", inFligth)
			time.Sleep(10 * time.Second)
			infoLogger.Printf("%d call pending\n", inFligth)
		} else if inFligth > 200 {
			infoLogger.Printf("%d call pending, waiting a little bit\n", inFligth)
			time.Sleep(10 * time.Second)
			infoLogger.Printf("%d call still pending\n", inFligth)
		}
	}
}
