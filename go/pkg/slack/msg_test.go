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

func TestParseSimpleCommand(t *testing.T) {
	type args struct {
		botID UserID
		msg   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Basic",
			args: args{
				botID: UserID{"12345673"},
				msg:   "<@12345673> roll me 2d20k1",
			},
			want: "roll me 2d20k1",
		},
		{
			name: "Not the bot",
			args: args{
				botID: UserID{"12345673"},
				msg:   "<@99999999> roll me 2d20k1",
			},
			want: "",
		},
		{
			name: "Bot not prefixed",
			args: args{
				botID: UserID{"12345673"},
				msg:   "roll me 2d20k1 <@12345673>",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSimpleCommand(tt.args.botID, tt.args.msg); got != tt.want {
				t.Errorf("ParseSimpleCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
