package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dymm/gorchestrator/pkg/config"
	"github.com/dymm/gorchestrator/pkg/data"
)

const outgoingQueue = "orchestrator"
const valueToAdd = 7

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()

	fmt.Println(("Starting processorAdd"))
	defer fmt.Println(("Stoping processorAdd"))

	for {
		workItem, err := myMessageQueue.Receive()

		var val data.TestValue
		if err == nil {
			val, err = data.DeserializeTestValue(workItem.GetValues())
		}
		if err != nil {
			fmt.Println("processorAdd : error while reading the message. ", err)
			os.Exit(0)
		}

		val.Value = val.Value + valueToAdd
		serializedValue, _ := json.Marshal(val)
		workItem.GetValues()["data"] = string(serializedValue)

		if val.Value >= 0 && val.Value <= 100 {
			fmt.Printf("%s processorAdd : loosing the value\n", val.Name)
			continue //Lose the message for a timeout
		}

		err = myMessageQueue.Send(outgoingQueue, workItem)
		if err != nil {
			fmt.Println("processorAdd : error while sending the message. ", err)
			os.Exit(0)
		}
	}

}
