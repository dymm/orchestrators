package workflow

import (
	"github.com/dymm/gorchestrator/pkg/messaging"
)

// Execute start the workflow
func Execute(theWorkflow Workflow, message messaging.WorkItem) error {

	var err error
	data := message.GetData()
	for _, step := range theWorkflow.Steps {
		data, err = step.Process(data)
		if err != nil {
			return err
		}
	}

	return nil
}
