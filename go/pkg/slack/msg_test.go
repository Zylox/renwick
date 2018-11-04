package slack

import (
	"reflect"
	"testing"
)

func TestDetectUsers(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want []UserID
	}{
		{
			name: "Basic",
			args: args{"<@UDSEXD9CD> other words"},
			want: []UserID{{"UDSEXD9CD"}},
		},
		{
			name: "Two Users",
			args: args{"<@UDSEXD9CD> other <@TACOTUESDAY> words"},
			want: []UserID{{"UDSEXD9CD"}, {"TACOTUESDAY"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectUsers(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
