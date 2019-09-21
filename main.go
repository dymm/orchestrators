package main

import (
	"fmt"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/messaging/localchannel"
	"github.com/dymm/gorchestrator/pkg/workflow"
)

func main() {
	allProducerAndConsumerQueues := getAllQueueOrDie()
	allWorflows := getTheWorkflowsOrDie(allProducerAndConsumerQueues)

	startTheProcessorsAndProducer(allProducerAndConsumerQueues)
	myMessageQueue := allProducerAndConsumerQueues["orchestrator"]

	for {
		workItem, err := myMessageQueue.Receive()
		var selectedWorkflow workflow.Workflow
		var workflowInfo workflow.Information
		if err == nil {
			selectedWorkflow, workflowInfo, err = workflow.SelectWorkflow(allWorflows, workItem)
		}

		var finished bool
		if err == nil {
			finished, err = workflow.SendToTheProcessor(selectedWorkflow, workflowInfo, workItem)
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

func getTheWorkflowsOrDie(queues map[string]messaging.Queue) []workflow.Workflow {

	return []workflow.Workflow{
		workflow.New("Value lower than 50", returnTrueIfTheValueIsLowerThan50,
			[]workflow.Step{
				workflow.Step{
					Name:    "Step 1",
					Process: queues["addConstToValue"],
				},
				workflow.Step{
					Name:    "Step 2",
					Process: queues["printTheValue"],
				},
			},
		),
		workflow.New("Value greater or equal than 50", returnTrueIfTheValueIsGreaterOrEqualThan50,
			[]workflow.Step{
				workflow.Step{
					Name:    "Step 1",
					Process: queues["subConstToValue"],
				},
				workflow.Step{
					Name:    "Step 2",
					Process: queues["subConstToValue"],
				},
				workflow.Step{
					Name:    "Step 3",
					Process: queues["printTheValue"],
				},
			},
		),
	}
}

func getAllQueueOrDie() map[string]messaging.Queue {
	queues := make(map[string]messaging.Queue)
	queues["orchestrator"] = localchannel.New()
	queues["addConstToValue"] = localchannel.New()
	queues["subConstToValue"] = localchannel.New()
	queues["printTheValue"] = localchannel.New()
	return queues
}

func startTheProcessorsAndProducer(queues map[string]messaging.Queue) {
	orchestratorQueue := queues["orchestrator"]
	go addConstToValue(queues["addConstToValue"], orchestratorQueue)
	go subConstToValue(queues["subConstToValue"], orchestratorQueue)
	go printTheValue(queues["printTheValue"], orchestratorQueue)
	go createValueProducer(orchestratorQueue)
}
