package main

import (
	"fmt"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/workflow"
)

func main() {

	allProducerAndConsumerQueues := getAllQueueOrDie()
	startTheProcessorsAndProducer(allProducerAndConsumerQueues)
	myMessageQueue := allProducerAndConsumerQueues["orchestrator"]

	allWorflows := getTheWorkflowsOrDie()

	for {
		workItem, err := myMessageQueue.Receive()
		var selectedWorkflow workflow.Workflow
		var workflowInfo workflow.Information
		if err == nil {
			selectedWorkflow, workflowInfo, err = workflow.SelectWorkflow(allWorflows, workItem)
		}

		var finished bool
		if err == nil {
			finished, err = workflow.SendToTheProcessor(myMessageQueue, selectedWorkflow, workflowInfo, workItem)
		}
		if err != nil {
			fmt.Println("Error while executing the workflow", err)
			return
		}
		if finished {
			fmt.Println("Workflow finished")
		}
	}
}

func getTheWorkflowsOrDie() []workflow.Workflow {

	return []workflow.Workflow{
		workflow.New("Value lower than 50",
			workflow.Validator{Value: "data.Value", Regex: `^(\d|[0-5]\d?)$`}, //50 or less
			[]workflow.Step{
				workflow.Step{
					Name:    "Step 1",
					Process: "addConstToValue",
				},
				workflow.Step{
					Name:    "Step 2",
					Process: "printTheValue",
				},
			},
		),
		workflow.New("Value greater or equal than 50",
			workflow.Validator{Value: "data.Value", Regex: `^([6-9]\d|\d{3,})$`}, //Greater than 50,
			[]workflow.Step{
				workflow.Step{
					Name:    "Step 1",
					Process: "subConstToValue",
				},
				workflow.Step{
					Name:    "Step 2",
					Process: "subConstToValue",
				},
				workflow.Step{
					Name:    "Step 3",
					Process: "printTheValue",
				},
			},
		),
	}
}

func getAllQueueOrDie() map[string]messaging.Queue {
	queues := make(map[string]messaging.Queue)
	queues["orchestrator"] = createMessageQueueOrDie("orchestrator")
	queues["addConstToValue"] = createMessageQueueOrDie("addConstToValue")
	queues["subConstToValue"] = createMessageQueueOrDie("subConstToValue")
	queues["printTheValue"] = createMessageQueueOrDie("printTheValue")
	queues["producer"] = createMessageQueueOrDie("producer")
	return queues
}

func startTheProcessorsAndProducer(queues map[string]messaging.Queue) {
	orchestratorQ := "orchestrator"
	go addConstToValue(queues["addConstToValue"], orchestratorQ)
	go subConstToValue(queues["subConstToValue"], orchestratorQ)
	go printTheValue(queues["printTheValue"], orchestratorQ)
	go createValueProducer(queues["producer"], orchestratorQ)
}
