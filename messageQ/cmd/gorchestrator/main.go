package main

import (
	"fmt"
	"time"

	"github.com/dymm/orchestrators/messageQ/pkg/config"
	"github.com/dymm/orchestrators/messageQ/pkg/workflow"
)

func main() {
	myMessageQueue := config.CreateMQMessageQueueOrDie()

	allWorflows := getTheWorkflowsOrDie()

	workflow.StartSessionTimeoutChecking(myMessageQueue, myMessageQueue.GetName())

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
					Process:   "processor_sub",
					OnSuccess: "Step 2",
					OnError:   "Dump",
					Timeout:   2,
				},
				"Step 2": workflow.Step{
					Process: "processor_print",
				},
				"Dump": workflow.Step{
					Process: "processor_error",
				},
			},
		),
		workflow.New("Value greater or equal than 50",
			workflow.Validator{Value: "data.Value", Regex: `^([6-9]\d|\d{3,})$`}, //Greater than 50,
			"Step 1",
			map[string]workflow.Step{
				"Step 1": workflow.Step{
					Process:   "processor_add",
					OnSuccess: "Step 2",
					OnError:   "Dump",
					Timeout:   2,
				},
				"Step 2": workflow.Step{
					Process:   "processor_add",
					OnSuccess: "Step 3",
					OnError:   "Dump",
					Timeout:   2,
				},
				"Step 3": workflow.Step{
					Process: "processor_print",
				},
				"Dump": workflow.Step{
					Process: "processor_error",
				},
			},
		),
	}
}
