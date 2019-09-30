package workflow

import (
	"reflect"
	"testing"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

func Test_getInformationFromWorkItem(t *testing.T) {
	type args struct {
		workItem messaging.WorkItem
	}
	tests := []struct {
		name string
		args args
		want Information
	}{
		{
			name: "Test with all needed informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `{"AssignedWorkflow":1, "CurrentStep":2}`})},
			want: Information{AssignedWorkflow: 1, CurrentStep: 2},
		},
		{
			name: "Test with missing informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `{"AssignedWorkflow":1`})},
			want: Information{AssignedWorkflow: -1, CurrentStep: -1},
		},
		{
			name: "Test with false informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `{"AssignedWorkflow":-10, "CurrentStep":2`})},
			want: Information{AssignedWorkflow: -1, CurrentStep: -1},
		},
		{
			name: "Test without informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `??`})},
			want: Information{AssignedWorkflow: -1, CurrentStep: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInformationFromWorkItem(tt.args.workItem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInformationFromWorkItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
