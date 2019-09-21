package main

import (
	"fmt"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/messaging/localchannel"
	"github.com/dymm/gorchestrator/pkg/workflow"
)

func main() {

	myWorkflow := getTheWorkflowOrDie()
	myMessageQueue := getTheMessageQueueOrDie()
	startTheProducer(myMessageQueue)

	for {
		message, err := myMessageQueue.Receive()
		if err == nil {
			err = workflow.Execute(myWorkflow, message)
		}
		if err != nil {
			fmt.Println("Error while executing the workflow", err)
			return
		}
	}
}

func getTheWorkflowOrDie() workflow.Workflow {

	return workflow.Workflow{
		Name: "test",
		Steps: []workflow.Step{
			workflow.Step{
				Name:    "Step 1",
				Process: addConstToValue,
			},
			workflow.Step{
				Name:    "Step 2",
				Process: printTheValue,
			},
		},
	}
}

func getTheMessageQueueOrDie() messaging.Queue {
	return localchannel.New()
}

func startTheProducer(queue messaging.Queue) {
	go createValue(queue)
}
