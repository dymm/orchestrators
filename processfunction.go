package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

type dataType struct {
	Name  string
	Value int
}

func createValue(queue messaging.Queue) {
	counter := 0
	for {
		counter = counter + 1
		newValue := dataType{
			Name:  fmt.Sprintf("Value %d", counter),
			Value: rand.Intn(100),
		}
		fmt.Printf("%s : Producing the value %d\n", newValue.Name, newValue.Value)
		if err := queue.Send(messaging.NewWorkItem(newValue)); err != nil {
			fmt.Println("Error while sending the message. ", err)
			os.Exit(0)
		}
		time.Sleep(3 * time.Second)
	}
}

func addConstToValue(input interface{}) (interface{}, error) {
	data, ok := input.(dataType)
	if !ok {
		return nil, errors.New("Can't cast the input data to the rigth type")
	}

	fmt.Printf("%s : Adding 1\n", data.Name)
	data.Value = data.Value + 1
	return data, nil
}

func printTheValue(input interface{}) (interface{}, error) {
	data, ok := input.(dataType)
	if !ok {
		return nil, errors.New("Can't cast the input data to the rigth type")
	}

	fmt.Printf("%s : value is %d\n", data.Name, data.Value)
	data.Value = data.Value + 1
	return data, nil
}
