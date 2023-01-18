package router

// endpoint is an address definition
type endpoint struct {
	addr    string
	healthy bool
}

func (e *endpoint) MarkUnhealthy() {
	e.healthy = false
}
