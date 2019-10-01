package workflow

import (
	"encoding/json"
	"errors"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//New return a new workflow
func New(name string, validator Validator, firstStep string, steps map[string]Step) Workflow {
	return Workflow{
		Name:      name,
		validate:  validator,
		FirstStep: firstStep,
		Steps:     steps,
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

	//Get the current step if the workflow is already running
	var currentStep Step
	currentFound := true
	if info.CurrentStep != "" {
		currentStep, currentFound = theWorkflow.Steps[info.CurrentStep]

	} else {
		currentStep.OnSuccess = theWorkflow.FirstStep //Set the first step for the workflow
	}
	if currentFound == false {
		return true, errors.New("No step '" + info.CurrentStep + "' found")
	}

	//Get the step name to execute now
	var nextStepName string
	if _, errorPresent := workItem.GetValues()["error"]; errorPresent {
		delete(workItem.GetValues(), "error") //Do not forward the error to the next step
		nextStepName = currentStep.OnError
	} else {
		nextStepName = currentStep.OnSuccess
	}

	if nextStepName == "" {
		return true, nil //No next step defined, the workflow is finished
	}

	//Get the step to execute now
	nextStep, nextFound := theWorkflow.Steps[nextStepName]
	if nextFound == false {
		return true, errors.New("No step '" + nextStepName + "' found")
	}

	//Send to the process
	info.CurrentStep = nextStepName
	serializedWorkflowInfo, _ := json.Marshal(info)
	values := workItem.GetValues()
	values["workflow"] = string(serializedWorkflowInfo)
	return false, queue.Send(nextStep.Process, messaging.NewWorkItem(values))
}
