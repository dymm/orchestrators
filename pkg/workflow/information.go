package workflow

import (
	"encoding/json"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//Information about the assigned workflow
type Information struct {
	AssignedWorkflow int
	CurrentStep      string
}

//getInformationFromWorkItem return the workflow information from a work item
func getInformationFromWorkItem(workItem messaging.WorkItem) Information {

	var info Information
	err := json.Unmarshal([]byte(workItem.GetValues()["workflow"]), &info)
	if err != nil {
		//No workflow information, it's a new incomming message
		info = Information{AssignedWorkflow: -1, CurrentStep: ""}
	}
	if err == nil && (info.AssignedWorkflow < 0 || info.CurrentStep == "") {
		info = Information{AssignedWorkflow: -1, CurrentStep: ""}
	}
	return info
}
