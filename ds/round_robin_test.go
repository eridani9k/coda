package ds

import (
	"testing"
)

func TestStep(t *testing.T) {
	tests := map[string]struct {
		s     []string
		index int
		want  int
	}{
		"increment_nonlooping_01": {
			s:     []string{"a", "b", "c", "d"},
			index: 0,
			want:  1,
		},
		"increment_nonlooping_02": {
			s:     []string{"a", "b", "c", "d"},
			index: 1,
			want:  2,
		},
		"increment_looping": {
			s:     []string{"a", "b", "c", "d"},
			index: 3,
			want:  0,
		},
		"single_element": {
			s:     []string{"a"},
			index: 0,
			want:  0,
		},
	}

	for name, ts := range tests {
		t.Run(name, func(t *testing.T) {
			got := step(ts.s, ts.index)
			if got != ts.want {
				t.Errorf("got: %+v, want: %+v", got, ts.want)
			}
		})
	}
}
