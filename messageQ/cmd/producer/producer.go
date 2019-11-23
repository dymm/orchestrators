package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/dymm/orchestrators/messageQ/pkg/config"
	"github.com/dymm/orchestrators/messageQ/pkg/data"
	"github.com/dymm/orchestrators/messageQ/pkg/messaging"
)

const outgoingQueue = "orchestrator"
const timeLayout = "2006-01-02 15:04:05.000000"

var inFligth int64
var infoLogger *log.Logger

func main() {
	infoLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	myMessageQueue := config.CreateMQMessageQueueOrDie()
	go receiveLoop(myMessageQueue)

	time.Sleep(3 * time.Second)
	var counter int64
	for {
		counter = counter + 1
		newValue := data.TestValue{
			Name:  fmt.Sprintf("Value %d", counter),
			Value: int(counter % 200),
		}

		current := time.Now()
		serialized, _ := json.Marshal(newValue)

		newWorkItem := messaging.NewWorkItem(
			map[string]string{
				"data":         string(serialized),
				"id":           strconv.FormatInt(counter, 10),
				"start":        current.Format(timeLayout),
				"finalReplyTo": myMessageQueue.GetName(),
			})

		if err := myMessageQueue.Send(outgoingQueue, newWorkItem); err != nil {
			infoLogger.Println("Error while sending the message. ", err)
			os.Exit(0)
		}
		atomic.AddInt64(&inFligth, 1)
		if counter%222 == 0 {
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

func receiveLoop(queue messaging.Queue) {

	for {
		workItem, err := queue.Receive()
		if err != nil {
			infoLogger.Println("Error while receiving a message", err)
			os.Exit(-1)
		}
		id := workItem.GetValues()["id"]
		startTime, err := time.Parse(timeLayout, workItem.GetValues()["start"])
		elapsed := time.Since(startTime)
		infoLogger.Printf("%s - End of processing in %s\n", id, elapsed)

		atomic.AddInt64(&inFligth, -1)
	}
}
