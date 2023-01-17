package ds

import (
	"fmt"
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

func (r *RoundRobin) Size() {
	return r.size
}
