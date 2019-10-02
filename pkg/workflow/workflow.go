package workflow

import (
	"errors"
	"strconv"

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

// GetTheWorkflowAndSession return the workflow that can handle the message
func GetTheWorkflowAndSession(allWorkflows []Workflow, workItem messaging.WorkItem) (Workflow, *Session, error) {

	var workflow Workflow

	session, found := getSessionStoredInTheWorkItem(workItem)
	if found == true {
		workflow = allWorkflows[session.assignedWorkflow]
	} else {
		session = createNewSession()
		for index, oneWorkflow := range allWorkflows {
			if oneWorkflow.CanHandleTheMessage(workItem.GetValues()) {
				workflow = oneWorkflow
				session.assignedWorkflow = index
				break
			}
		}
	}

	if len(workflow.Name) == 0 {
		return workflow, session, errors.New("No workflow found for the data")
	}
	return workflow, session, nil
}

// SendToTheProcessor send the data to the processor
// return true if the workflow is finished and an error if needed
func SendToTheProcessor(queue messaging.Queue, theWorkflow Workflow, session *Session, workItem messaging.WorkItem) (bool, error) {

	//Get the current step if the workflow is already running
	var currentStep Step
	currentFound := true
	if session.currentStep != "" {
		currentStep, currentFound = theWorkflow.Steps[session.currentStep]

	} else {
		currentStep.OnSuccess = theWorkflow.FirstStep //Set the first step for the workflow
	}
	if currentFound == false {
		return true, errors.New("No step '" + session.currentStep + "' found")
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
	session.currentStep = nextStepName
	values := workItem.GetValues()
	values["sessionId"] = strconv.FormatUint(session.Key, 10)
	return false, queue.Send(nextStep.Process, messaging.NewWorkItem(values))
}
