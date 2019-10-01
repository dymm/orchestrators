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
			fmt.Printf("Workflow '%d' finished\n", workflowInfo.AssignedWorkflow)
		}
	}
}

func getTheWorkflowsOrDie() []workflow.Workflow {

	return []workflow.Workflow{
		workflow.New("Value lower than 50",
			workflow.Validator{Value: "data.Value", Regex: `^(\d|[0-5]\d?)$`}, //50 or less
			"Step 1",
			map[string]workflow.Step{
				"Step 1": workflow.Step{
					Process:   "subConstToValue",
					OnSuccess: "Step 2",
					OnError:   "Dump",
				},
				"Step 2": workflow.Step{
					Process: "printTheValue",
				},
				"Step 3": workflow.Step{
					Process: "dumpTheValue",
				},
				"Dump": workflow.Step{
					Process: "dumpTheValue",
				},
			},
		),
		workflow.New("Value greater or equal than 50",
			workflow.Validator{Value: "data.Value", Regex: `^([6-9]\d|\d{3,})$`}, //Greater than 50,
			"Step 1",
			map[string]workflow.Step{
				"Step 1": workflow.Step{
					Process:   "addConstToValue",
					OnSuccess: "Step 2",
					OnError:   "Dump",
				},
				"Step 2": workflow.Step{
					Process:   "addConstToValue",
					OnSuccess: "Step 3",
					OnError:   "Dump",
				},
				"Step 3": workflow.Step{
					Process: "printTheValue",
				},
				"Dump": workflow.Step{
					Process: "dumpTheValue",
				},
			},
		),
	}
}

func getAllQueueOrDie() map[string]messaging.Queue {
	queues := make(map[string]messaging.Queue)
	queues["orchestrator"] = createMessageQueueOrDie("orchestrator", queues)
	queues["addConstToValue"] = createMessageQueueOrDie("addConstToValue", queues)
	queues["subConstToValue"] = createMessageQueueOrDie("subConstToValue", queues)
	queues["printTheValue"] = createMessageQueueOrDie("printTheValue", queues)
	queues["dumpTheValue"] = createMessageQueueOrDie("dumpTheValue", queues)
	queues["producer"] = createMessageQueueOrDie("producer", queues)
	return queues
}

func startTheProcessorsAndProducer(queues map[string]messaging.Queue) {
	orchestratorQ := "orchestrator"
	go addConstToValue(queues["addConstToValue"], orchestratorQ, 23)
	go subConstToValue(queues["subConstToValue"], orchestratorQ, 14)
	go printTheValue(queues["printTheValue"], orchestratorQ)
	go dumpTheValue(queues["dumpTheValue"], orchestratorQ)
	go createValueProducer(queues["producer"], orchestratorQ)
}
