package router

// endpoint is an address definition of a server
// fronted by the Router.
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

func (e *endpoint) setUnhealthy() {
	e.healthy = false
}

func (e *endpoint) setHealthy() {
	e.healthy = true
}
