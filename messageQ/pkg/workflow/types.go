package workflow

import "regexp"

//Validator is a couple of value name and a regex
type Validator struct {
	Value string
	Regex string
}

// Workflow is a workflow description
type Workflow struct {
	//Name of the workflow
	Name     string
	validate Validator
	//The first step too be executed
	FirstStep string
	//Steps to be executed
	Steps map[string]Step
}

//CanHandleTheMessage return true if the workflow is siuted for the data
func (w Workflow) CanHandleTheMessage(values map[string]string) bool {

	value, err := getStringFromJSONMap(w.validate.Value, values)
	if err != nil {
		return false
	}
	var regex *regexp.Regexp
	regex, err = regexp.Compile(w.validate.Regex)
	if err != nil {
		return false
	}
	return regex.MatchString(value)
}

//Step is a processing to apply to the datas
type Step struct {
	//Process is a message queue name to a processor
	Process string
	//Execute the specified step on success
	OnSuccess string
	//Execute the specified step on error
	OnError string
	//Timeout in second
	Timeout uint
}
