package router

import (
	"reflect"
	"testing"
)

var (
	emptyRouter = &Router{}

	singleEndpointRouter = &Router{
		endpoints: []endpoint{
			{addr: "8080", valid: true},
		},
		curr: 0,
		next: 0,
		size: 1,
	}

	multipleEndpointRouter = &Router{
		endpoints: []endpoint{
			{addr: "8080", valid: true},
			{addr: "8081", valid: true},
			{addr: "8082", valid: true},
			{addr: "8083", valid: true},
		},
		curr: 0,
		next: 1,
		size: 4,
	}

	RouterWithInvalidEndpoints = &Router{
		endpoints: []endpoint{
			{addr: "8080", valid: true},
			{addr: "8081", valid: false},
			{addr: "8082", valid: true},
			{addr: "8083", valid: false},
			{addr: "8084", valid: false},
			{addr: "8085", valid: true},
		},
		curr: 0,
		next: 1,
		size: 6,
	}
)

func TestNewRouter(t *testing.T) {
	tests := map[string]struct {
		addrs []string
		want  *Router
	}{
		"empty_addr": {
			addrs: []string{},
			want:  nil,
		},
		"single_addr": {
			addrs: []string{"8080"},
			want:  singleEndpointRouter,
		},
		"multiple_addrs": {
			addrs: []string{"8080", "8081", "8082", "8083"},
			want:  multipleEndpointRouter,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewRouter(ts.addrs)
			if !reflect.DeepEqual(got, ts.want) {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestSeekValid(t *testing.T) {
	tests := map[string]struct {
		router *Router
		index  int
		want   int
	}{
		"t1": {
			router: RouterWithInvalidEndpoints,
			index:  0,
			want:   2,
		},
		"t2": {
			router: RouterWithInvalidEndpoints,
			index:  2,
			want:   5,
		},
		"t3": {
			router: RouterWithInvalidEndpoints,
			index:  5,
			want:   0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ts.router.SeekValid(ts.index)
			if err != nil {
				panic(err)
			}

			if got != ts.want {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestStep(t *testing.T) {
	tests := map[string]struct {
		router *Router
		index  int
		want   int
	}{
		"step_nonlooping": {
			router: multipleEndpointRouter,
			index:  1,
			want:   2,
		},
		"step_looping": {
			router: multipleEndpointRouter,
			index:  3,
			want:   0,
		},
		"step_single_element": {
			router: singleEndpointRouter,
			index:  0,
			want:   0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := ts.router.Step(ts.index)
			if got != ts.want {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestUtilStep(t *testing.T) {
	tests := map[string]struct {
		max   int
		index int
		want  int
	}{
		"increment_nonlooping_01": {
			max:   3,
			index: 0,
			want:  1,
		},
		"increment_nonlooping_02": {
			max:   3,
			index: 1,
			want:  2,
		},
		"increment_looping": {
			max:   3,
			index: 3,
			want:  0,
		},
		"single_element": {
			max:   0,
			index: 0,
			want:  0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := step(ts.max, ts.index)
			if got != ts.want {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}
