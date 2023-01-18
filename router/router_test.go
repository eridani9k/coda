package router

import (
	"reflect"
	"testing"
)

var (
	emptyRouter = &Router{}

	singleEndpointRouter = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
		},
		curr: 0,
		next: 0,
		size: 1,
	}

	multipleEndpointRouter = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8081", healthy: true},
			{addr: "8082", healthy: true},
			{addr: "8083", healthy: true},
		},
		curr: 0,
		next: 1,
		size: 4,
	}

	RouterWithUnhealthyEndpoints = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8081", healthy: false},
			{addr: "8082", healthy: true},
			{addr: "8083", healthy: false},
			{addr: "8084", healthy: false},
			{addr: "8085", healthy: true},
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
			want:  emptyRouter,
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

func TestSeekHealthy(t *testing.T) {
	tests := map[string]struct {
		router *Router
		index  int
		want   int
		err    error
	}{
		"t1": {
			router: RouterWithUnhealthyEndpoints,
			index:  0,
			want:   2,
			err:    nil,
		},
		"t2": {
			router: RouterWithUnhealthyEndpoints,
			index:  2,
			want:   5,
			err:    nil,
		},
		"t3": {
			router: RouterWithUnhealthyEndpoints,
			index:  5,
			want:   0,
			err:    nil,
		},
		"t4": {
			router: emptyRouter,
			index:  1,
			want:   0,
			err:    ErrNoValidEndpointFound,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ts.router.SeekHealthy(ts.index)
			if got != ts.want || err != ts.err {
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
		"step_single_element": {
			router: singleEndpointRouter,
			index:  0,
			want:   0,
		},
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
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := ts.router.step(ts.index)
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
		"single_element": {
			max:   0,
			index: 0,
			want:  0,
		},
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
