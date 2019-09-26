package workflow

import (
	"encoding/json"
	"errors"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//New return a new workflow
func New(name string, validator Validator, steps []Step) Workflow {
	return Workflow{
		Name:     name,
		validate: validator,
		Steps:    steps,
	}
}

// SelectWorkflow return the workflow that can handle the message
func SelectWorkflow(allWorkflows []Workflow, workItem messaging.WorkItem) (Workflow, Information, error) {

	info := getInformationFromWorkItem(workItem)

	var workflow Workflow
	if info.AssignedWorkflow >= 0 && info.AssignedWorkflow < len(allWorkflows) {
		workflow = allWorkflows[info.AssignedWorkflow]
	} else {
		for index, oneWorkflow := range allWorkflows {
			if oneWorkflow.CanHandleTheMessage(workItem.GetValues()) {
				workflow = oneWorkflow
				info.AssignedWorkflow = index
				break
			}
		}
	}
	if len(workflow.Name) == 0 {
		return workflow, info, errors.New("No workflow found for the data")
	}
	return workflow, info, nil
}

// SendToTheProcessor send the data to the processor
// return true if the workflow is finished and an error if needed
func SendToTheProcessor(queue messaging.Queue, theWorkflow Workflow, info Information, workItem messaging.WorkItem) (bool, error) {
	info.CurrentStep = info.CurrentStep + 1
	if info.CurrentStep >= len(theWorkflow.Steps) {
		return true, nil
	}

	serializedWorkflowInfo, _ := json.Marshal(info)
	values := workItem.GetValues()
	values["workflow"] = string(serializedWorkflowInfo)
	return false, queue.Send(theWorkflow.Steps[info.CurrentStep].Process, messaging.NewWorkItem(values))
}
