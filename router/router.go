package router

import (
	"errors"
)

var (
	ErrNoEndpointsRegistered = errors.New("no endpoints registered")
)

type Router struct {
	endpoints []*endpoint
	curr      int
	next      int
	size      int
}

func NewRouter(addrs []string) *Router {
	if len(addrs) == 0 {
		return &Router{}
	}

	curr := 0
	next := step(len(addrs)-1, curr)

	endpoints := make([]*endpoint, len(addrs))
	for i, addr := range addrs {
		endpoints[i] = &endpoint{
			addr:    addr,
			healthy: true,
		}
	}

	return &Router{
		endpoints: endpoints,
		curr:      curr,
		next:      next,
		size:      len(addrs),
	}
}

// NextEndpoint returns the endpoint at r.curr and advances
// both r.curr and r.next to their next valid positions.
func (r *Router) NextEndpoint() (*endpoint, error) {
	if r.NoEndpoints() {
		return nil, ErrNoEndpointsRegistered
	}

	endpoint := r.getEndpoint()
	r.curr = r.next
	r.next = r.seekHealthy(r.next)

	return endpoint, nil
}

func (r *Router) getEndpoint() *endpoint {
	return r.endpoints[r.curr]
}

// seekHealthy returns the index of the next healthy
// endpoint starting from the element after index.
func (r *Router) seekHealthy(index int) int {
	k := r.step(index)
	for k != index {
		if r.endpoints[k].healthy {
			return k
		}
		k = r.step(k)
	}

	// All other endpoints are unhealthy or index is the only element.
	return index
}

// step is only called when len(r.endpoints) > 0.
// Length checks are done on the ancestor function.
func (r *Router) step(index int) int {
	return step(len(r.endpoints)-1, index)
}

// max is 0-indexed.
// TODO: Convert argument to uint?
func step(max int, index int) int {
	if index == max {
		return 0
	}

	return index + 1
}

func (r *Router) Register(addr string) {
	r.endpoints = append(r.endpoints, &endpoint{
		addr:    addr,
		healthy: true,
	})
	r.size += 1
}

func (r *Router) Size() int {
	return r.size
}

func (r *Router) NoEndpoints() bool {
	return r.size == 0
}
