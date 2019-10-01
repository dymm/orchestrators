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
			name: "Two JSON object but the latest is wrong",
			args: args{
				value:  "first.second.name",
				values: map[string]string{"first": `{"second":{error}}`},
			},
			want:    "",
			wantErr: true,
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
		{
			name: "Two JSON object on the same level but searching for an unkown value",
			args: args{
				value:  "first.second.value",
				values: map[string]string{"first": `{"other":{"name":"good"},"second":{"name":"found"}}`},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Three JSON object",
			args: args{
				value:  "first.second.third.name",
				values: map[string]string{"first": `{"other":{"name":"good"},"second":{"third":{"name":"i'm the third"}}}`},
			},
			want:    "i'm the third",
			wantErr: false,
		},
		{
			name: "Three JSON object but reading a wrong value",
			args: args{
				value:  "first.second.third.name",
				values: map[string]string{"first": `{"other":{"name":"good"},"second":{ $$, third":{"name":"i'm the third"}}}`},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getStringFromJSONMap(tt.args.value, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStringFromJSONMap() '%s' error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getStringFromJSONMap() '%s' = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
