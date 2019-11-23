package main

import (
	"fmt"
	"os"

	"github.com/dymm/orchestrators/messageQ/pkg/config"
	"github.com/dymm/orchestrators/messageQ/pkg/data"
)

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()

	fmt.Println(("Starting processor-print"))
	defer fmt.Println(("Stoping processor-print"))

	for {
		workItem, err := myMessageQueue.Receive()

		var val data.TestValue
		if err == nil {
			val, err = data.DeserializeTestValue(workItem.GetValues())
		}
		if err != nil {
			fmt.Println("processor-print : error while reading the message. ", err)
			os.Exit(0)
		}

		fmt.Printf("%s : value is %d\n", val.Name, val.Value)

		err = myMessageQueue.Send(workItem.GetValues()["replyTo"], workItem)
		if err != nil {
			fmt.Println("processor-print : error while sending the message. ", err)
			os.Exit(0)
		}
	}

}
