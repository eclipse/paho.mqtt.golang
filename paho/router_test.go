package paho

import (
	"reflect"
	"testing"
)

func Test_match(t *testing.T) {
	tests := []struct {
		name  string
		route string
		topic string
		want  bool
	}{
		{"basic1", "a/b", "a/b", true},
		{"basic2", "a", "a/b", false},
		{"plus1", "a/+", "a/b", true},
		{"plus2", "+/b", "a/b", true},
		{"plus3", "a/+/c", "a/b/c", true},
		{"plus4", "a/+/c", "a/asdf/c", true},
		{"hash1", "#", "a/b", true},
		{"hash2", "a/#", "a/b", true},
		{"hash3", "b/#", "a/b", false},
		{"hash4", "#", "", true},
		{"share1", "$share/a/b", "a/b", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := match(tt.route, tt.topic); got != tt.want {
				t.Errorf("match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routeIncludesTopic(t *testing.T) {
	type args struct {
		route string
		topic string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := routeIncludesTopic(tt.args.route, tt.args.topic); got != tt.want {
				t.Errorf("routeIncludesTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routeSplit(t *testing.T) {
	type args struct {
		route string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := routeSplit(tt.args.route); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("routeSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}
