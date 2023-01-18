package ds

import (
	"reflect"
	"testing"
)

var (
	emptyRR = &RoundRobin{}

	singleEndpointRR = &RoundRobin{
		endpoints: []endpoint{
			{addr: "8080", valid: true},
		},
		curr: 0,
		next: 0,
		size: 1,
	}

	nonEmptyRR = &RoundRobin{
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
)

func TestNewRR(t *testing.T) {
	tests := map[string]struct {
		addrs []string
		want  *RoundRobin
	}{
		"empty_addr": {
			addrs: []string{},
			want:  nil,
		},
		"single_addr": {
			addrs: []string{"8080"},
			want:  singleEndpointRR,
		},
		"multiple_addrs": {
			addrs: []string{"8080", "8081", "8082", "8083"},
			want:  nonEmptyRR,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewRoundRobin(ts.addrs)
			if !reflect.DeepEqual(got, ts.want) {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestRRStep(t *testing.T) {
	tests := map[string]struct {
		rr    *RoundRobin
		index int
		want  int
	}{
		"step_nonlooping": {
			rr:    nonEmptyRR,
			index: 1,
			want:  2,
		},
		"step_looping": {
			rr:    nonEmptyRR,
			index: 3,
			want:  0,
		},
		"step_single_element": {
			rr:    singleEndpointRR,
			index: 0,
			want:  0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := ts.rr.Step(ts.index)
			if got != ts.want {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}

func TestStep(t *testing.T) {
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
