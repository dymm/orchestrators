package workflow

import "testing"

func Test_getStringFromJSONMap(t *testing.T) {
	type args struct {
		value  string
		values map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Single value",
			args: args{
				value:  "first",
				values: map[string]string{"first": "found"},
			},
			want:    "found",
			wantErr: false,
		},
		{
			name: "Unknown value",
			args: args{
				value:  "?",
				values: map[string]string{"first": "found"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "One JSON object",
			args: args{
				value:  "first.name",
				values: map[string]string{"first": `{"name":"found"}`},
			},
			want:    "found",
			wantErr: false,
		},
		{
			name: "Two JSON object",
			args: args{
				value:  "first.second.name",
				values: map[string]string{"first": `{"second":{"name":"found"}}`},
			},
			want:    "found",
			wantErr: false,
		},
		{
			name: "Two JSON object on the same level",
			args: args{
				value:  "first.second.name",
				values: map[string]string{"first": `{"other":{"name":"good"},"second":{"name":"found"}}`},
			},
			want:    "found",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStringFromJSONMap(tt.args.value, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStringFromJSONMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStringFromJSONMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
