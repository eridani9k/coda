package ds

import (
// "fmt"
)

type RoundRobin struct {
	addrs []string
	curr  int
	next  int
	size  int
}

func NewRoundRobin(addrs []string) *RoundRobin {
	if addrs == nil || len(addrs) == 0 {
		return nil
	}

	var curr int
	var next int
	if len(addrs) == 1 {
		next = curr
	} else {
		next = curr + 1
	}

	return &RoundRobin{
		addrs: addrs,
		curr:  curr,
		next:  next,
		size:  len(addrs),
	}
}

func (r *RoundRobin) Add(addr string) {
	r.addrs = append(r.addrs, addr)
	r.size += 1
}

func (r *RoundRobin) Size() int {
	return r.size
}

func (r *RoundRobin) Step(index int) int {
	return step(len(r.addrs)-1, index)
}

// step is 0-indexed.
func step(max int, index int) int {
	if index == max {
		return 0
	}

	return index + 1
}
