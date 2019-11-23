package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dymm/orchestrators/messageQ/pkg/config"
	"github.com/dymm/orchestrators/messageQ/pkg/data"
)

const valueToSub = 7

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()

	fmt.Println(("Starting processor-sub"))
	defer fmt.Println(("Stoping processor-sub"))

	for {
		workItem, err := myMessageQueue.Receive()

		var val data.TestValue
		if err == nil {
			val, err = data.DeserializeTestValue(workItem.GetValues())
		}
		if err != nil {
			fmt.Println("processor-sub : error while reading the message. ", err)
			os.Exit(0)
		}

		val.Value = val.Value + valueToSub
		serializedValue, _ := json.Marshal(val)
		workItem.GetValues()["data"] = string(serializedValue)

		if val.Value >= 45 && val.Value <= 55 {
			fmt.Printf("%s processor-sub : loosing the value\n", val.Name)
			continue //Lose the message for a timeout
		}

		err = myMessageQueue.Send(workItem.GetValues()["replyTo"], workItem)
		if err != nil {
			fmt.Println("processor-sub : error while sending the message. ", err)
			os.Exit(0)
		}
	}

}
