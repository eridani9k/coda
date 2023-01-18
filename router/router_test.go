package router

import (
	"reflect"
	"testing"
)

// Refer to test_vars.go for structures used in this file.

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
			want:  routerWithSingleEndpoint,
		},
		"multiple_addrs": {
			addrs: []string{"8080", "8081", "8082", "8083"},
			want:  routerWithMultipleEndpoints,
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
		"multiple_unhealthy_01": {
			router: routerWithUnhealthyEndpointsV1,
			index:  0,
			want:   2,
			err:    nil,
		},
		"multiple_unhealthy_02": {
			router: routerWithUnhealthyEndpointsV1,
			index:  2,
			want:   5,
			err:    nil,
		},
		"multiple_unhealthy_03": {
			router: routerWithUnhealthyEndpointsV1,
			index:  5,
			want:   0,
			err:    nil,
		},
		"multiple_unhealthy_04": {
			router: routerWithUnhealthyEndpointsV2,
			index:  2,
			want:   4,
			err:    nil,
		},
		"multiple_unhealthy_05": {
			router: routerWithUnhealthyEndpointsV2,
			index:  4,
			want:   2,
			err:    nil,
		},
		"single_healthy_01": {
			router: routerWithSingleHealthyEndpoint,
			index:  3,
			want:   3,
			err:    nil,
		},
		"single_healthy_02": {
			router: routerWithSingleEndpoint,
			index:  0,
			want:   0,
			err:    nil,
		},
		"emptyRouter": {
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
			router: routerWithSingleEndpoint,
			index:  0,
			want:   0,
		},
		"step_nonlooping": {
			router: routerWithMultipleEndpoints,
			index:  1,
			want:   2,
		},
		"step_looping": {
			router: routerWithMultipleEndpoints,
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
