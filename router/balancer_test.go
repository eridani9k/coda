package router

import (
	"reflect"
	"testing"
)

// Refer to test_vars.go for structures used in this file.

func TestNewBalancer(t *testing.T) {
	tests := map[string]struct {
		addrs []string
		want  *Balancer
	}{
		"empty_addr": {
			addrs: []string{},
			want:  emptyBalancer,
		},
		"single_addr": {
			addrs: []string{":8080"},
			want:  balancerWithSingleEndpoint,
		},
		"multiple_addrs": {
			addrs: []string{":8080", ":8081", ":8082", ":8083"},
			want:  balancerWithMultipleEndpoints,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewBalancer(ts.addrs)
			if !reflect.DeepEqual(got, ts.want) {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestAdvance(t *testing.T) {
	tests := map[string]struct {
		balancer *Balancer
		endpoint *endpoint
		newCurr  int
		newNext  int
		err      error
	}{
		"empty": {
			balancer: emptyBalancer,
			endpoint: nil,
			newCurr:  -1,
			err:      ErrNoValidEndpoints,
		},
		"single_endpoint": {
			balancer: balancerWithSingleEndpoint,
			endpoint: &endpoint{
				addr:    ":8080",
				healthy: true,
			},
			newCurr: 0,
			err:     nil,
		},
		"multiple_endpoints": {
			balancer: balancerWithMultipleEndpoints,
			endpoint: &endpoint{
				addr:    ":8081",
				healthy: true,
			},
			newCurr: 1,
			err:     nil,
		},
		"double_endpoints_01": {
			balancer: balancerAdvanceV1,
			endpoint: &endpoint{
				addr:    ":8080",
				healthy: true,
			},
			newCurr: 0,
			err:     nil,
		},
		"double_endpoints_02": {
			balancer: balancerAdvanceV2,
			endpoint: &endpoint{
				addr:    ":8081",
				healthy: true,
			},
			newCurr: 1,
			err:     nil,
		},
		"single_healthy_endpoint": {
			balancer: balancerAdvanceV3,
			endpoint: &endpoint{
				addr:    ":8083",
				healthy: true,
			},
			newCurr: 3,
			err:     nil,
		},
		"mixed_endpoints_01": {
			balancer: balancerAdvanceV4,
			endpoint: &endpoint{
				addr:    ":8081",
				healthy: true,
			},
			newCurr: 1,
			err:     nil,
		},
		"mixed_endpoints_05": {
			balancer: balancerAdvanceV5,
			endpoint: &endpoint{
				addr:    ":8080",
				healthy: true,
			},
			newCurr: 0,
			err:     nil,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			endpoint, err := ts.balancer.Advance()
			if err != ts.err || !reflect.DeepEqual(endpoint, ts.endpoint) || ts.balancer.curr != ts.newCurr {
				t.Errorf("\nendpoint - got: %+v, want: %+v\ncurr - got: %+v, want: %+v\n", endpoint, ts.endpoint, ts.balancer.curr, ts.newCurr)
			}
		})
	}
}

func TestSeekHealthy(t *testing.T) {
	tests := map[string]struct {
		balancer *Balancer
		index    int
		want     int
	}{
		"all_unhealthy": {
			balancer: balancerAllUnhealthy,
			index:    3,
			want:     -1,
		},
		"multiple_unhealthy_01": {
			balancer: balancerWithUnhealthyEndpointsV1,
			index:    0,
			want:     2,
		},
		"multiple_unhealthy_02": {
			balancer: balancerWithUnhealthyEndpointsV1,
			index:    2,
			want:     5,
		},
		"multiple_unhealthy_03": {
			balancer: balancerWithUnhealthyEndpointsV1,
			index:    5,
			want:     0,
		},
		"multiple_unhealthy_04": {
			balancer: balancerWithUnhealthyEndpointsV2,
			index:    2,
			want:     4,
		},
		"multiple_unhealthy_05": {
			balancer: balancerWithUnhealthyEndpointsV2,
			index:    4,
			want:     2,
		},
		"single_healthy_01": {
			balancer: balancerWithSingleHealthyEndpoint,
			index:    3,
			want:     3,
		},
		"single_healthy_02": {
			balancer: balancerWithSingleEndpoint,
			index:    0,
			want:     0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := ts.balancer.seekHealthy(ts.index)
			if got != ts.want {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestStep(t *testing.T) {
	tests := map[string]struct {
		balancer *Balancer
		index    int
		want     int
	}{
		"step_single_element": {
			balancer: balancerWithSingleEndpoint,
			index:    0,
			want:     0,
		},
		"step_nonlooping": {
			balancer: balancerWithMultipleEndpoints,
			index:    1,
			want:     2,
		},
		"step_looping": {
			balancer: balancerWithMultipleEndpoints,
			index:    3,
			want:     0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := ts.balancer.step(ts.index)
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
