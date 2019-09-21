package workflow

import (
	"errors"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//New return a new workflow
func New(name string, validateFunc Validator, steps []Step) Workflow {
	return Workflow{
		Name:     name,
		validate: validateFunc,
		Steps:    steps,
	}
}

// SelectWorkflow return the workflow that can handle the message
func SelectWorkflow(allWorkflows []Workflow, message messaging.WorkItem) (Workflow, error) {

	var workflow Workflow
	for _, oneWorkflow := range allWorkflows {
		if oneWorkflow.CanHandleTheMessage(message.GetData()) {
			workflow = oneWorkflow
			break
		}
	}
	if len(workflow.Name) == 0 {
		return workflow, errors.New("No workflow found for the data")
	}
	return workflow, nil
}

// Execute start the workflow
func Execute(theWorkflow Workflow, message messaging.WorkItem) error {

	var err error
	data := message.GetData()
	for _, step := range theWorkflow.Steps {
		data, err = step.Process(data)
		if err != nil {
			return err
		}
	}

	return nil
}
