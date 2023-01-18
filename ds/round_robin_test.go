package ds

import (
	"testing"
)

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
