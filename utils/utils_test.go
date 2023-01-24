package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTimestampMsg(t *testing.T) {
	ts := time.Now()
	msg := "Ping failed for :8080"

	got := timestampMsg(msg, ts)
	want := fmt.Sprintf("[ %s ] %s", ts.Format(time.RFC3339), msg)

	if got != want {
		t.Errorf("got: %+v, want: %+v", got, want)
	}
}
