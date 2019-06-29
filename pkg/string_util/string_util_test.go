package string_util

import (
	"testing"
)

func TestFirstLetterToLowerCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty string", args: args{str: ""}, want: ""},
		{name: "one letter upper string", args: args{str: "A"}, want: "a"},
		{name: "one letter lower string", args: args{str: "a"}, want: "a"},
		{name: "long upper string", args: args{str: "AAAAA"}, want: "aAAAA"},
		{name: "long lower string", args: args{str: "aaaaa"}, want: "aaaaa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstLetterToLowerCase(tt.args.str); got != tt.want {
				t.Errorf("FirstLetterToLowerCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstLetterToUpperCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty string", args: args{str: ""}, want: ""},
		{name: "one letter lower string", args: args{str: "a"}, want: "A"},
		{name: "one letter upper string", args: args{str: "A"}, want: "A"},
		{name: "long lower string", args: args{str: "aaaaa"}, want: "Aaaaa"},
		{name: "long upper string", args: args{str: "AAAAA"}, want: "AAAAA"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstLetterToUpperCase(tt.args.str); got != tt.want {
				t.Errorf("FirstLetterToUpperCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
