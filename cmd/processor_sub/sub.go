package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dymm/gorchestrator/pkg/config"
	"github.com/dymm/gorchestrator/pkg/data"
)

const outgoingQueue = "orchestrator"
const valueToSub = 7

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()

	fmt.Println(("Starting processor_sub"))
	defer fmt.Println(("Stoping processor_sub"))

	for {
		workItem, err := myMessageQueue.Receive()

		var val data.TestValue
		if err == nil {
			val, err = data.DeserializeTestValue(workItem.GetValues())
		}
		if err != nil {
			fmt.Println("processor_sub : error while reading the message. ", err)
			os.Exit(0)
		}

		val.Value = val.Value + valueToSub
		serializedValue, _ := json.Marshal(val)
		workItem.GetValues()["data"] = string(serializedValue)

		if val.Value >= 35 && val.Value <= 65 {
			fmt.Printf("%s processor_sub : loosing the value\n", val.Name)
			continue //Lose the message for a timeout
		}

		err = myMessageQueue.Send(outgoingQueue, workItem)
		if err != nil {
			fmt.Println("processor_sub : error while sending the message. ", err)
			os.Exit(0)
		}
	}

}
