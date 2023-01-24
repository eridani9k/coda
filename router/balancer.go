package router

import (
	"errors"
	"sync"
)

var (
	ErrNoValidEndpoints = errors.New("no valid endpoints")
)

type Balancer struct {
	endpoints []*endpoint
	curr      int // -1 indicates no valid endpoints.
	size      int

	mu sync.Mutex
}

func NewBalancer(addrs []string) *Balancer {
	var curr int
	if len(addrs) == 0 {
		curr = -1
	}

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
		size:      len(addrs),
	}
}

// Advance() is the mechanism for moving b.curr and retrieving an endpoint.
// Advance() seeks for the next healthy endpoint, moves b.curr
// to that index, and returns the *endpoint at the new index.
// Advance() exits early if there are no valid endpoints (b.curr == -1).
func (b *Balancer) Advance() (*endpoint, error) {
	if b.curr == -1 {
		return nil, ErrNoValidEndpoints
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	err := b.nextHealthy()
	if err != nil {
		return nil, err
	}

	return b.getEndpoint(), nil
}

func (b *Balancer) getEndpoint() *endpoint {
	return b.endpoints[b.curr]
}

// nextHealthy() moves b.curr to the next healthy endpoint.
// nextHealthy() assumes pre-processed b.curr is >= 0.
// If there are no valid endpoints found, sets b.curr to -1
// and returns an error.
func (b *Balancer) nextHealthy() error {
	b.curr = b.seekHealthy(b.curr)
	if b.curr == -1 {
		return ErrNoValidEndpoints
	}

	return nil
}

// seekHealthy() returns the index of the next healthy
// endpoint starting from the endpoint after index.
// seekHealthy() returns -1 if there is no healthy endpoint found.
// seekHealthy() will check for health at most len(endpoints) times.
// TODO: rewrite this description.
func (b *Balancer) seekHealthy(index int) int {
	k := b.step(index)
	for k != index {
		if b.endpoints[k].healthy {
			return k
		}
		k = b.step(k)
	}

	// k has looped back to index.
	// If the endpoint at index is unhealthy, return -1.
	if !b.endpoints[index].healthy {
		return -1
	}

	return index
}

// step() returns the index of the next element.
// If index has reached the end of the collection,
// loop back to 0.
func (b *Balancer) step(index int) int {
	return step(len(b.endpoints)-1, index)
}

// step() is a helper function that returns the next
// integer after index. If index exceeds max, loop back to 0.
// step() is the foundational mechanism for Balancer's
// round-robin traversal.
func step(max int, index int) int {
	if index == max {
		return 0
	}

	return index + 1
}

// Add() adds an endpoint instance to b.endpoints.
// New endpoints are always marked by default as healthy.
func (b *Balancer) Add(addr string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.endpoints = append(b.endpoints, &endpoint{
		addr:    addr,
		healthy: true,
	})
	b.size += 1

	// Reset b.curr to 0 if there were no valid endpoints previously.
	if b.curr == -1 {
		b.curr = 0
	}
}

func (b *Balancer) ResetCurr() {
	if b.size == 0 {
		return
	}

	b.curr = 0
}

func (b *Balancer) Curr() int {
	return b.curr
}

func (b *Balancer) Size() int {
	return b.size
}
