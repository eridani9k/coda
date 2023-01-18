package router

import (
	"errors"
	// "fmt"
)

var (
	ErrNoValidEndpointFound = errors.New("no valid endpoint found")
)

// An endpoint is an address definition
type endpoint struct {
	addr    string
	healthy bool
}

type Router struct {
	endpoints []endpoint
	curr      int
	next      int
	size      int
}

func NewRouter(addrs []string) *Router {
	if len(addrs) == 0 {
		return nil
	}

	curr := 0
	next := step(len(addrs)-1, curr)

	endpoints := make([]endpoint, len(addrs))
	for i, addr := range addrs {
		endpoints[i] = endpoint{
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

func (r *Router) Add(addr string) {
	r.endpoints = append(r.endpoints, endpoint{
		addr:    addr,
		healthy: true,
	})
	r.size += 1
}

func (r *Router) Size() int {
	return r.size
}

// SeekValid returns the index of the next valid address from element after index.
func (r *Router) SeekHealthy(index int) int {
	k := r.Step(index)
	for k != index {
		if r.endpoints[k].healthy {
			return k
		}
		k = r.Step(k)
	}

	// All other endpoints are unhealthy, return index.
	return index
}

func (r *Router) Step(index int) int {
	return step(r.size-1, index)
}

// max is 0-indexed.
func step(max int, index int) int {
	if index == max {
		return 0
	}

	return index + 1
}

func (r *Router) MarkUnhealthy(index int) {
	r.endpoints[index].healthy = false
}
