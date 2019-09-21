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
func SelectWorkflow(allWorkflows []Workflow, info Information) (Workflow, error) {
	var workflow Workflow
	if info.assignedWorkflow >= 0 && info.assignedWorkflow < len(allWorkflows) {
		workflow = allWorkflows[info.assignedWorkflow]
	} else {
		for index, oneWorkflow := range allWorkflows {
			if oneWorkflow.CanHandleTheMessage(info.GetData()) {
				workflow = oneWorkflow
				info.assignedWorkflow = index
				break
			}
		}
	}
	if len(workflow.Name) == 0 {
		return workflow, errors.New("No workflow found for the data")
	}
	return workflow, nil
}

// SendToTheProcessor send the data to the processor
// return true if the workflow is finished and an error if needed
func SendToTheProcessor(theWorkflow Workflow, info Information) (bool, error) {
	info.currentStep = info.currentStep + 1
	if info.currentStep >= len(theWorkflow.Steps) {
		return true, nil
	}
	return false, theWorkflow.Steps[info.currentStep].Process.Send(messaging.NewWorkItem(info))
}
