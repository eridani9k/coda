package router

import (
	"fmt"
)

// endpoint is an address definition
type endpoint struct {
	addr    string
	healthy bool
}

func (e *endpoint) getAddress() string {
	return e.addr
}

func (e *endpoint) isHealthy() bool {
	return e.healthy
}

func (e *endpoint) markUnhealthy() {
	fmt.Printf("Endpoint %s marked unhealthy, and will not receive requests.\n", e.addr)
	e.healthy = false
}
