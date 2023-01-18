package ds

import (
//"fmt"
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
	return step(r.addrs, index)
}

// TODO: abstract this further to step(max int, index int).
// Passing the slice is unnecessary.
func step(s []string, index int) int {
	if index == len(s)-1 {
		return 0
	}

	return index + 1
}
