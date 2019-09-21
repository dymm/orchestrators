package main

import (
	"fmt"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/messaging/localchannel"
	"github.com/dymm/gorchestrator/pkg/workflow"
)

func main() {

	allWorflows := getTheWorkflowsOrDie()
	myMessageQueue := getTheMessageQueueOrDie()
	startTheProducer(myMessageQueue)

	for {
		var selectedWorkflow workflow.Workflow
		message, err := myMessageQueue.Receive()
		if err == nil {
			selectedWorkflow, err = workflow.SelectWorkflow(allWorflows, message)
		}
		if err == nil {
			err = workflow.Execute(selectedWorkflow, message)
		}
		if err != nil {
			fmt.Println("Error while executing the workflow", err)
			return
		}
	}
}

func getTheWorkflowsOrDie() []workflow.Workflow {

	return []workflow.Workflow{
		workflow.New("Value lower than 50", returnTrueIfTheValueIsLowerThan50,
			[]workflow.Step{
				workflow.Step{
					Name:    "Step 1",
					Process: addConstToValue,
				},
				workflow.Step{
					Name:    "Step 2",
					Process: printTheValue,
				},
			},
		),
		workflow.New("Value greater or equal than 50", returnTrueIfTheValueIsGreaterOrEqualThan50,
			[]workflow.Step{
				workflow.Step{
					Name:    "Step 1",
					Process: subConstToValue,
				},
				workflow.Step{
					Name:    "Step 2",
					Process: subConstToValue,
				},
				workflow.Step{
					Name:    "Step 3",
					Process: printTheValue,
				},
			},
		),
	}
}

func getTheMessageQueueOrDie() messaging.Queue {
	return localchannel.New()
}

func startTheProducer(queue messaging.Queue) {
	go createValue(queue)
}
