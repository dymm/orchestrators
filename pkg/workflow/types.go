package workflow

// Workflow is a workflow description
type Workflow struct {
	//Name of the workflow
	Name string
	//Steps to be executed
	Steps []Step
}

//Step is a processing to apply to the datas
type Step struct {
	//Name of the step
	Name string
	//Process is a function that do something to the data given in parameter and returing the result
	Process func(interface{}) (interface{}, error)
}
