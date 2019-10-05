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

// SendToTheNextProcessor send the data to the processor
// return true if the workflow is finished and an error if needed
func SendToTheNextProcessor(queue messaging.Queue, theWorkflow Workflow, session *Session, workItem messaging.WorkItem) (bool, error) {

	finished, step, workItem, err := getTheNextStep(theWorkflow, session, workItem)
	if finished == false && err == nil {
		setStepInformationInSession(session, step, workItem)
		err = queue.Send(step.Process, workItem)
	}
	//Send to the process
	return finished, err
}

//GetTheNextStep return true if the workflow is finished, if not return a workitem to send and update the session
// If a next step is found the session variable 'CurrentStep' is updated
func getTheNextStep(theWorkflow Workflow, session *Session, workItem messaging.WorkItem) (bool, Step, messaging.WorkItem, error) {
	//Get the current step if the workflow is already running
	var currentStep Step
	currentFound := true
	if session.CurrentStep.Name != "" {
		currentStep, currentFound = theWorkflow.Steps[session.CurrentStep.Name]
	} else {
		currentStep.OnSuccess = theWorkflow.FirstStep //Set the first step for the workflow
	}
	if currentFound == false {
		return true, Step{}, workItem, errors.New("No step '" + session.CurrentStep.Name + "' found")
	}

	//Get the step name to execute now
	if _, errorPresent := workItem.GetValues()["error"]; errorPresent {
		session.CurrentStep.Name = currentStep.OnError
	} else {
		session.CurrentStep.Name = currentStep.OnSuccess
	}
	if session.CurrentStep.Name == "" {
		return true, Step{}, workItem, nil //No next step defined, the workflow is finished
	}

	//Get the step to execute now
	nextStep, nextFound := theWorkflow.Steps[session.CurrentStep.Name]
	if nextFound == false {
		return true, Step{}, workItem, errors.New("No step '" + session.CurrentStep.Name + "' found")
	}

	values := workItem.GetValues()
	values["sessionId"] = strconv.FormatUint(session.Key, 10)
	newWorkItem := messaging.NewWorkItem(values)

	return false, nextStep, newWorkItem, nil
}
