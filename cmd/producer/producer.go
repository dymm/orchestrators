package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dymm/gorchestrator/pkg/config"
	"github.com/dymm/gorchestrator/pkg/data"
	"github.com/dymm/gorchestrator/pkg/messaging"
)

const outgoingQueue = "orchestrator"

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()
	time.Sleep(3 * time.Second)
	counter := 0
	for {
		counter = counter + 1
		newValue := data.TestValue{
			Name:  fmt.Sprintf("Value %d", counter),
			Value: rand.Intn(100),
		}
		fmt.Printf("%s : Producing the value %d\n", newValue.Name, newValue.Value)
		serialized, _ := json.Marshal(newValue)
		newWorkItem := messaging.NewWorkItem(map[string]string{"data": string(serialized)})

		if err := myMessageQueue.Send(outgoingQueue, newWorkItem); err != nil {
			fmt.Println("Error while sending the message. ", err)
			os.Exit(0)
		}
		if counter%222 == 0 {
			time.Sleep(30 * time.Second)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
