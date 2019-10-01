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
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `{"AssignedWorkflow":1, "CurrentStep":"second"}`})},
			want: Information{AssignedWorkflow: 1, CurrentStep: "second"},
		},
		{
			name: "Test with missing informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `{"AssignedWorkflow":1`})},
			want: Information{AssignedWorkflow: -1, CurrentStep: ""},
		},
		{
			name: "Test with false informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `{"AssignedWorkflow":-10, "CurrentStep":"second"`})},
			want: Information{AssignedWorkflow: -1, CurrentStep: ""},
		},
		{
			name: "Test without informations",
			args: args{workItem: messaging.NewWorkItem(map[string]string{"workflow": `??`})},
			want: Information{AssignedWorkflow: -1, CurrentStep: ""},
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
