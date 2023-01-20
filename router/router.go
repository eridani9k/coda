package router

import (
	"errors"
	"sync"
)

var (
	ErrNoEndpointsRegistered = errors.New("no endpoints registered")
)

type Router struct {
	endpoints []*endpoint
	curr      int
	next      int
	size      int

	mu sync.Mutex
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

// Advance returns the endpoint at r.curr and advances
// both r.curr and r.next to their next valid positions.
func (r *Router) Advance() (*endpoint, error) {
	if r.NoEndpoints() {
		return nil, ErrNoEndpointsRegistered
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	endpoint := r.getEndpoint()
	r.curr = r.next
	r.next = r.seekHealthy(r.next)

	return endpoint, nil
}

// Peek returns the endpoint at r.next, but does not
// advance r.curr or r.next.
func (r *Router) Peek() (*endpoint, error) {
	if r.NoEndpoints() {
		return nil, ErrNoEndpointsRegistered
	}

	return r.endpoints[r.next], nil
}

func (r *Router) getEndpoint() *endpoint {
	return r.endpoints[r.curr]
}

// seekHealthy returns the index of the next healthy
// endpoint starting from the endpoint after index.
// seekHealthy assumes the endpoint at index is healthy,
// and returns index if:
//   - r.endpoints only contain the endpoint at index,
//   - all other endpoints in r.endpoints are unhealthy.
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

// step returns the index of the next element.
// If index has reached the end of the collection,
// loop back to 0.
func (r *Router) step(index int) int {
	return step(len(r.endpoints)-1, index)
}

// step returns the next integer after index.
// If index exceeds max, loop back to 0.
func step(max int, index int) int {
	if index == max {
		return 0
	}

	return index + 1
}

// Register adds an endpoint instance to r.endpoints.
// New endpoints are always marked by default as healthy.
func (r *Router) Register(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
