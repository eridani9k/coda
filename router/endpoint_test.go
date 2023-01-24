package router

import (
	"reflect"
	"testing"
)

func TestGetAddress(t *testing.T) {
	e := &endpoint{
		addr:    "http://127.0.0.1:10000",
		healthy: true,
	}

	got := e.getAddress()
	want := "http://127.0.0.1:10000"

	if got != want {
		t.Errorf("got: %+v, want: %+v", got, want)
	}
}

func TestIsHealthy(t *testing.T) {
	e := &endpoint{
		addr:    "http://127.0.0.1:10000",
		healthy: true,
	}

	got := e.isHealthy()
	want := true

	if got != want {
		t.Errorf("got: %+v, want: %+v", got, want)
	}
}

func TestSetHealthy(t *testing.T) {
	e := &endpoint{
		addr:    "http://127.0.0.1:10000",
		healthy: false,
	}

	e.setHealthy()
	got := e
	want := &endpoint{
		addr:    "http://127.0.0.1:10000",
		healthy: true,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %+v, want: %+v", got, want)
	}
}

func TestSetUnhealthy(t *testing.T) {
	e := &endpoint{
		addr:    "http://127.0.0.1:10000",
		healthy: true,
	}

	e.setUnhealthy()
	got := e
	want := &endpoint{
		addr:    "http://127.0.0.1:10000",
		healthy: false,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %+v, want: %+v", got, want)
	}
}
