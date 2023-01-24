package router

// endpoint is an address definition of a server
// fronted by the Router.
type endpoint struct {

	// addr is a fully qualified address of a backend
	// server. This value mirrors that found in the
	// addresses.cfg file.
	addr string

	// healthy indicates if this endpoint is ready to
	// receive requests. This value is modified when ever
	// a change in connectivity is detected, either through
	// a health check, or a failed proxy request.
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
