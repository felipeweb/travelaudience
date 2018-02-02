package numbers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/felipeweb/travelaudience/fakeapi"
)

func Test_removeDuplicates(t *testing.T) {
	type args struct {
		elements []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"without duplicates", args{[]int{1, 3, 2, 4}}, []int{1, 3, 2, 4}},
		{"with duplicates", args{[]int{1, 3, 3, 2, 2, 4}}, []int{1, 3, 2, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicates(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderHandler(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		url string
	}
	tests := []struct {
		name       string
		args       args
		empty      bool
		statusCode int
		hasParam   bool
	}{
		{"without urls", args{httptest.NewRecorder(), "/"}, false, http.StatusBadRequest, false},
		{"with urls and numbers result", args{httptest.NewRecorder(), "/?u=%s/fibo"}, false, http.StatusOK, true},
		{"with urls and empty result", args{httptest.NewRecorder(), "/?u=%s/fibo"}, true, http.StatusOK, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(fakeapi.Mux(tt.empty))
			defer server.Close()
			if tt.hasParam {
				tt.args.url = fmt.Sprintf(tt.args.url, server.URL)
			}
			OrderHandler(tt.args.w, httptest.NewRequest("GET", tt.args.url, nil))
			resp := tt.args.w.Result()
			if resp.StatusCode != tt.statusCode {
				t.Errorf("OrderHandler() = %v, but got %v", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_exec(t *testing.T) {
	type args struct {
		ctx  context.Context
		urls []string
	}
	tests := []struct {
		name  string
		args  args
		empty bool
		want  []int
	}{
		{"with numbers", args{context.Background(), []string{"%s/fibo"}}, false, []int{1, 2, 3, 5, 8, 13, 21}},
		{"empty", args{context.Background(), []string{"%s/fibo"}}, true, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(fakeapi.Mux(tt.empty))
			defer server.Close()
			for i := range tt.args.urls {
				tt.args.urls[i] = fmt.Sprintf(tt.args.urls[i], server.URL)
			}
			if got := exec(tt.args.ctx, tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
