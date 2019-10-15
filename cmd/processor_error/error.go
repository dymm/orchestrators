package main

import (
	"fmt"
	"os"

	"github.com/dymm/gorchestrator/pkg/config"
	"github.com/dymm/gorchestrator/pkg/data"
)

const outgoingQueue = "orchestrator"

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()

	fmt.Println(("Starting processorError"))
	defer fmt.Println(("Stoping processorError"))

	for {
		workItem, err := myMessageQueue.Receive()

		var val data.TestValue
		if err == nil {
			val, err = data.DeserializeTestValue(workItem.GetValues())
		}
		if err != nil {
			fmt.Println("processorError : error while reading the message. ", err)
			os.Exit(0)
		}

		fmt.Printf("Error on %s with a value %d\n", val.Name, val.Value)

		err = myMessageQueue.Send(outgoingQueue, workItem)
		if err != nil {
			fmt.Println("processorError : error while sending the message. ", err)
			os.Exit(0)
		}
	}

}
