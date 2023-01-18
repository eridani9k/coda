package router

// This file contains variables used in router_test.go.

var (
	emptyRouter = &Router{}

	routerWithSingleEndpoint = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
		},
		curr: 0,
		next: 0,
		size: 1,
	}

	routerWithMultipleEndpoints = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8081", healthy: true},
			{addr: "8082", healthy: true},
			{addr: "8083", healthy: true},
		},
		curr: 0,
		next: 1,
		size: 4,
	}

	routerWithUnhealthyEndpointsV1 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8081", healthy: false},
			{addr: "8082", healthy: true},
			{addr: "8083", healthy: false},
			{addr: "8084", healthy: false},
			{addr: "8085", healthy: true},
		},
		curr: 0,
		next: 1,
		size: 6,
	}

	routerWithUnhealthyEndpointsV2 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: false},
			{addr: "8081", healthy: false},
			{addr: "8082", healthy: true},
			{addr: "8083", healthy: false},
			{addr: "8084", healthy: true},
			{addr: "8085", healthy: false},
		},
		curr: 0,
		next: 1,
		size: 6,
	}

	routerWithSingleHealthyEndpoint = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: false},
			{addr: "8081", healthy: false},
			{addr: "8082", healthy: false},
			{addr: "8083", healthy: true},
			{addr: "8084", healthy: false},
			{addr: "8085", healthy: false},
		},
		curr: 0,
		next: 1,
		size: 6,
	}
)
