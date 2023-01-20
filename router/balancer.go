package router

import (
	"errors"
	"sync"
)

var (
	ErrNoEndpointsRegistered = errors.New("no endpoints registered")
)

type Balancer struct {
	endpoints []*endpoint
	curr      int
	next      int
	size      int

	mu sync.Mutex
}

func NewBalancer(addrs []string) *Balancer {
	if len(addrs) == 0 {
		return &Balancer{}
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

	return &Balancer{
		endpoints: endpoints,
		curr:      curr,
		next:      next,
		size:      len(addrs),
	}
}

// Advance returns the endpoint at b.curr and advances
// both b.curr and b.next to their next valid positions.
func (b *Balancer) Advance() (*endpoint, error) {
	if b.NoEndpoints() {
		return nil, ErrNoEndpointsRegistered
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	endpoint := b.getEndpoint()
	b.curr = b.next
	b.next = b.seekHealthy(b.next)

	return endpoint, nil
}

// Peek returns the endpoint at b.next, but does not
// advance b.curr or r.next.
func (b *Balancer) Peek() (*endpoint, error) {
	if b.NoEndpoints() {
		return nil, ErrNoEndpointsRegistered
	}

	return b.endpoints[b.next], nil
}

func (b *Balancer) getEndpoint() *endpoint {
	return b.endpoints[b.curr]
}

// seekHealthy returns the index of the next healthy
// endpoint starting from the endpoint after index.
// seekHealthy assumes the endpoint at index is healthy,
// and returns index if:
//   - r.endpoints only contain the endpoint at index,
//   - all other endpoints in r.endpoints are unhealthy.
//
// TODO: seekHealthy is missing a NoValidEndpoint signal.
// TODO: return -1 if there is no healthy endpoint found?
func (b *Balancer) seekHealthy(index int) int {
	k := b.step(index)
	for k != index {
		if b.endpoints[k].healthy {
			return k
		}
		k = b.step(k)
	}

	// All other endpoints are unhealthy or index is the only element.
	return index
}

// step returns the index of the next element.
// If index has reached the end of the collection,
// loop back to 0.
func (b *Balancer) step(index int) int {
	return step(len(b.endpoints)-1, index)
}

// step returns the next integer after index.
// If index exceeds max, loop back to 0.
func step(max int, index int) int {
	if index == max {
		return 0
	}

	return index + 1
}

// Register adds an endpoint instance to b.endpoints.
// New endpoints are always marked by default as healthy.
func (b *Balancer) Register(addr string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.endpoints = append(b.endpoints, &endpoint{
		addr:    addr,
		healthy: true,
	})
	b.size += 1
}

func (b *Balancer) Size() int {
	return b.size
}

func (b *Balancer) NoEndpoints() bool {
	return b.size == 0
}
