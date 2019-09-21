package workflow

import "github.com/dymm/gorchestrator/pkg/messaging"

//Validator is a function made for validate a workfow for a given data
type Validator func(values map[string]string) bool

// Workflow is a workflow description
type Workflow struct {
	//Name of the workflow
	Name     string
	validate Validator
	//Steps to be executed
	Steps []Step
}

//CanHandleTheMessage return true if the workflow is siuted for the data
func (w Workflow) CanHandleTheMessage(values map[string]string) bool {
	return w.validate(values)
}

//Step is a processing to apply to the datas
type Step struct {
	//Name of the step
	Name string
	//Process is a message queue communicating with a processor
	Process messaging.Queue
}
