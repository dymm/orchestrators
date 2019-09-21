package workflow

//Validator is a function made for validate a workfow for a given data
type Validator func(interface{}) bool

// Workflow is a workflow description
type Workflow struct {
	//Name of the workflow
	Name         string
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
	//Process is a function that do something to the data given in parameter and returing the result
	Process func(interface{}) (interface{}, error)
}
