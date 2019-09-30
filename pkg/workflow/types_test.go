package workflow

import "testing"

func TestWorkflow_CanHandleTheMessage(t *testing.T) {
	type args struct {
		values map[string]string
	}
	tests := []struct {
		name string
		w    Workflow
		args args
		want bool
	}{
		{
			name: "Refuse a bad regex",
			w:    Workflow{validate: Validator{Value: "data.value", Regex: `(((`}},
			args: args{map[string]string{"data": `{"value":"61"}`}},
			want: false,
		},
		{
			name: "Accept a value equal or grater than 60",
			w:    Workflow{validate: Validator{Value: "data.value", Regex: `^([6-9]\d|\d{3,})$`}},
			args: args{map[string]string{"data": `{"value":"61"}`}},
			want: true,
		},
		{
			name: "Refuse value not equal or grater than 60",
			w:    Workflow{validate: Validator{Value: "data.value", Regex: `^([6-9]\d|\d{3,})$`}},
			args: args{map[string]string{"data": `{"value":"2"}`}},
			want: false,
		},
		{
			name: "Refuse a dictionnary without the good variable",
			w:    Workflow{validate: Validator{Value: "data.value", Regex: `^([6-9]\d|\d{3,})$`}},
			args: args{map[string]string{"something": `{"value":"2"}`}},
			want: false,
		},
		{
			name: "Refuse a dictionnary without the good value",
			w:    Workflow{validate: Validator{Value: "data.value", Regex: `^([6-9]\d|\d{3,})$`}},
			args: args{map[string]string{"data": `{"something":"2"}`}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.CanHandleTheMessage(tt.args.values); got != tt.want {
				t.Errorf("Workflow.CanHandleTheMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
