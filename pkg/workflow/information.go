package workflow

import "github.com/dymm/gorchestrator/pkg/messaging"

//Information about the assigned workflow
type Information struct {
	assignedWorkflow int
	currentStep      int
	data             interface{}
}

//GetData return the data
func (info Information) GetData() interface{} {
	return info.data
}

//GetInformationFromWorkItem return the workflow information from a work item
func GetInformationFromWorkItem(workItem messaging.WorkItem) Information {
	info, ok := workItem.GetData().(Information)
	if !ok {
		//No workflow information, it's a new incomming message
		info = Information{assignedWorkflow: -1, currentStep: -1, data: workItem.GetData()}
	}
	return info
}

//CreateWorkItemResponse create a new WorkItem that can be sent as a response
func CreateWorkItemResponse(workItem messaging.WorkItem, newData interface{}) messaging.WorkItem {
	info := GetInformationFromWorkItem(workItem)
	info.data = newData
	return messaging.NewWorkItem(info)
}
