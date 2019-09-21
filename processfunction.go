package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/workflow"
)

type dataType struct {
	Name  string
	Value int
}

func returnTrueIfTheValueIsLowerThan50(input interface{}) bool {
	data, ok := input.(dataType)
	if !ok {
		return false
	}
	return data.Value < 50
}

func returnTrueIfTheValueIsGreaterOrEqualThan50(input interface{}) bool {
	data, ok := input.(dataType)
	if !ok {
		return false
	}
	return data.Value >= 50
}

func createValueProducer(outgoing messaging.Queue) {
	time.Sleep(3 * time.Second)
	counter := 0
	for {
		counter = counter + 1
		newValue := dataType{
			Name:  fmt.Sprintf("Value %d", counter),
			Value: rand.Intn(100),
		}
		fmt.Printf("%s : Producing the value %d\n", newValue.Name, newValue.Value)
		if err := outgoing.Send(messaging.NewWorkItem(newValue)); err != nil {
			fmt.Println("Error while sending the message. ", err)
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
}

func addConstToValue(incoming messaging.Queue, outgoing messaging.Queue) {
	fmt.Println(("Starting addConstToValue"))
	defer fmt.Println(("Stoping addConstToValue"))

	for {
		workItem, err := incoming.Receive()
		if err != nil {
			fmt.Println("addConstToValue : error while reading the message. ", err)
			os.Exit(0)
		}

		info := workflow.GetInformationFromWorkItem(workItem)
		data, ok := info.GetData().(dataType)
		if !ok {
			fmt.Println("addConstToValue : Can't cast the input data to the rigth type")
			os.Exit(0)
		}

		fmt.Printf("%s : Adding 1\n", data.Name)
		data.Value = data.Value + 1
		err = outgoing.Send(workflow.CreateWorkItemResponse(workItem, data))
		if err != nil {
			fmt.Println("addConstToValue : error while sending the message. ", err)
			os.Exit(0)
		}
	}
}

func subConstToValue(incoming messaging.Queue, outgoing messaging.Queue) {
	fmt.Println(("Starting subConstToValue"))
	defer fmt.Println(("Stoping subConstToValue"))

	for {
		workItem, err := incoming.Receive()
		if err != nil {
			fmt.Println("subConstToValue : error while reading the message. ", err)
			os.Exit(0)
		}

		info := workflow.GetInformationFromWorkItem(workItem)
		data, ok := info.GetData().(dataType)
		if !ok {
			fmt.Println("subConstToValue : Can't cast the input data to the rigth type")
			os.Exit(0)
		}

		fmt.Printf("%s : Substracting 9\n", data.Name)
		data.Value = data.Value - 9
		err = outgoing.Send(workflow.CreateWorkItemResponse(workItem, data))
		if err != nil {
			fmt.Println("subConstToValue : error while sending the message. ", err)
			os.Exit(0)
		}
	}
}

func printTheValue(incoming messaging.Queue, outgoing messaging.Queue) {
	fmt.Println(("Starting printTheValue"))
	defer fmt.Println(("Stoping printTheValue"))

	for {
		workItem, err := incoming.Receive()
		if err != nil {
			fmt.Println("printTheValue : error while reading the message. ", err)
			os.Exit(0)
		}

		info := workflow.GetInformationFromWorkItem(workItem)
		data, ok := info.GetData().(dataType)
		if !ok {
			fmt.Println("printTheValue : Can't cast the input data to the rigth type")
			os.Exit(0)
		}

		fmt.Printf("%s : value is %d\n", data.Name, data.Value)

		err = outgoing.Send(workflow.CreateWorkItemResponse(workItem, data))
		if err != nil {
			fmt.Println("printTheValue : error while sending the message. ", err)
			os.Exit(0)
		}
	}
}
