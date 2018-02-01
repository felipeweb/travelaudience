package fakeapi

import (
	"net/http"
	"testing"
)

func Test_Mux(t *testing.T) {
	tests := []struct {
		name  string
		want  http.Handler
		empty bool
	}{
		{"mux false", Mux(false), false},
		{"mux true", Mux(true), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mux(tt.empty); got == nil {
				t.Errorf("Mux() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name  string
		args  args
		want  func(http.ResponseWriter, *http.Request)
		empty bool
	}{
		{"handler", args{[]int{1, 2, 3, 4}}, handler([]int{1, 2, 3, 4}, false), false},
		{"handler empty", args{[]int{1, 2, 3, 4}}, handler([]int{1, 2, 3, 4}, true), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler(tt.args.numbers, tt.empty); got == nil {
				t.Errorf("handler() = <nil>, want not nil")
			}
		})
	}
}
