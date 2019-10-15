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

	fmt.Println(("Starting processorPrint"))
	defer fmt.Println(("Stoping processorPrint"))

	for {
		workItem, err := myMessageQueue.Receive()

		var val data.TestValue
		if err == nil {
			val, err = data.DeserializeTestValue(workItem.GetValues())
		}
		if err != nil {
			fmt.Println("processorPrint : error while reading the message. ", err)
			os.Exit(0)
		}

		fmt.Printf("%s : value is %d\n", val.Name, val.Value)

		err = myMessageQueue.Send(outgoingQueue, workItem)
		if err != nil {
			fmt.Println("processorPrint : error while sending the message. ", err)
			os.Exit(0)
		}
	}

}
