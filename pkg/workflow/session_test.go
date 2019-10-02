package workflow

import (
	"reflect"
	"testing"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

func Test_getSessionStoredInTheWorkItem(t *testing.T) {

	session1 := createNewSession()
	session1.assignedWorkflow = 1
	session1.currentStep = "first"

	session2 := createNewSession()
	session2.assignedWorkflow = 1
	session2.currentStep = "second"

	type args struct {
		workItem messaging.WorkItem
	}
	tests := []struct {
		name        string
		args        args
		wantSession *Session
		wantFound   bool
	}{
		{
			name:        "Test with all needed informations",
			args:        args{workItem: messaging.NewWorkItem(map[string]string{"sessionId": `2`})},
			wantSession: session2,
			wantFound:   true,
		},
		{
			name:        "Test with all needed informations again",
			args:        args{workItem: messaging.NewWorkItem(map[string]string{"sessionId": `1`})},
			wantSession: session1,
			wantFound:   true,
		},
		{
			name:        "Test with missing informations",
			args:        args{workItem: messaging.NewWorkItem(map[string]string{"sessionId": ``})},
			wantSession: nil,
			wantFound:   false,
		},
		{
			name:        "Test with false informations",
			args:        args{workItem: messaging.NewWorkItem(map[string]string{"sessionId": `10`})},
			wantSession: nil,
			wantFound:   false,
		},
		{
			name:        "Test without informations",
			args:        args{workItem: messaging.NewWorkItem(map[string]string{"sessionId": `??`})},
			wantSession: nil,
			wantFound:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSession, gotFound := getSessionStoredInTheWorkItem(tt.args.workItem)
			if !reflect.DeepEqual(gotSession, tt.wantSession) {
				t.Errorf("getSessionStoredInTheWorkItem() %s gotSession = %v, want %v", tt.name, gotSession, tt.wantSession)
			}
			if gotFound != tt.wantFound {
				t.Errorf("getSessionStoredInTheWorkItem() %s gotFound = %v, want %v", tt.name, gotFound, tt.wantFound)
			}
		})
	}
}
