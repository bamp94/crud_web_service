package common

import (
	"reflect"
	"testing"
)

func TestFindUrlPathObjects(t *testing.T) {
	type args struct {
		urlPath string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{name: "result-1", args: args{"/api/v1/user/1/comment/4"}, want: [][]string{[]string{"api/v1", "user/1", "comment/4"}}},
		{name: "result-2", args: args{"/user/1/comment/4"}, want: [][]string{[]string{"user/1", "comment/4"}}},
		{name: "result-3", args: args{"/user/1"}, want: [][]string{[]string{"user/1"}}},
		{name: "result-4", args: args{"/user/"}, want: [][]string{[]string{"user/"}}},
		{name: "result-5", args: args{""}, want: [][]string{[]string{""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindUrlPathObjects(tt.args.urlPath); reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailToAscii() = %v, want %v", got, tt.want)
			}
		})
	}
}
