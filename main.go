package main

import (
	"fmt"
	"time"

	"github.com/dymm/gorchestrator/pkg/messaging"
	"github.com/dymm/gorchestrator/pkg/workflow"
)

func main() {

	allProducerAndConsumerQueues := getAllQueueOrDie()
	startTheProcessorsAndProducer(allProducerAndConsumerQueues)
	myMessageQueue := allProducerAndConsumerQueues["orchestrator"]

	allWorflows := getTheWorkflowsOrDie()

	workflow.StartSessionTimeoutChecking(myMessageQueue, "orchestrator")

	for {
		workItem, err := myMessageQueue.Receive()
		var selectedWorkflow workflow.Workflow
		var workflowSession *workflow.Session
		if err == nil {
			selectedWorkflow, workflowSession, err = workflow.GetTheWorkflowAndSession(allWorflows, workItem)
		}

		var finished bool
		if err == nil {
			finished, err = workflow.SendToTheNextProcessor(myMessageQueue, selectedWorkflow, workflowSession, workItem)
		}
		if err != nil {
			fmt.Printf("Error while executing the workflow %d. %s\n", workflowSession.Key, err)
			finished = true
		}
		if finished {
			fmt.Printf("Workflow '%d' finished in %d ms\n", workflowSession.Key, time.Now().Sub(workflowSession.CurrentStep.Started).Milliseconds())
			workflow.DeleteSession(workflowSession)
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
					Timeout:   2,
				},
				"Step 2": workflow.Step{
					Process: "printTheValue",
				},
				"Dump": workflow.Step{
					Process: "handleError",
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
					Timeout:   2,
				},
				"Step 2": workflow.Step{
					Process:   "addConstToValue",
					OnSuccess: "Step 3",
					OnError:   "Dump",
					Timeout:   2,
				},
				"Step 3": workflow.Step{
					Process: "printTheValue",
				},
				"Dump": workflow.Step{
					Process: "handleError",
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
	queues["handleError"] = createMessageQueueOrDie("handleError", queues)
	queues["producer"] = createMessageQueueOrDie("producer", queues)
	return queues
}

func startTheProcessorsAndProducer(queues map[string]messaging.Queue) {
	orchestratorQ := "orchestrator"
	for i := 0; i < 3; i++ {
		go addConstToValue(queues["addConstToValue"], orchestratorQ, 10)
		go subConstToValue(queues["subConstToValue"], orchestratorQ, 14)
		go printTheValue(queues["printTheValue"], orchestratorQ)
		go handleError(queues["handleError"], orchestratorQ)
	}
	go createValueProducer(queues["producer"], orchestratorQ)
}
