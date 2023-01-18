package router

// endpoint is an address definition
type endpoint struct {
	addr    string
	healthy bool
}

func (e *endpoint) GetAddress() string {
	return e.addr
}

func (e *endpoint) IsHealthy() bool {
	return e.healthy
}

func (e *endpoint) MarkUnhealthy() {
	e.healthy = false
}
