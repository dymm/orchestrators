package workflow

import "github.com/dymm/gorchestrator/pkg/messaging"

//Validator is a function made for validate a workfow for a given data
type Validator func(interface{}) bool

// Workflow is a workflow description
type Workflow struct {
	//Name of the workflow
	Name     string
	validate Validator
	//Steps to be executed
	Steps []Step
}

//CanHandleTheMessage return true if the workflow is siuted for the data
func (w Workflow) CanHandleTheMessage(data interface{}) bool {
	return w.validate(data)
}

//Step is a processing to apply to the datas
type Step struct {
	//Name of the step
	Name string
	//Process is a message queue communicating with a processor
	Process messaging.Queue
}
