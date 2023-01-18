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

	routerNextV1 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8081", healthy: true},
		},
		curr: 1,
		next: 0,
		size: 2,
	}

	routerNextV2 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8081", healthy: true},
		},
		curr: 0,
		next: 1,
		size: 2,
	}

	routerNextV3 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: false},
			{addr: "8081", healthy: false},
			{addr: "8082", healthy: false},
			{addr: "8083", healthy: true},
			{addr: "8084", healthy: false},
			{addr: "8085", healthy: false},
		},
		curr: 3,
		next: 3,
		size: 6,
	}

	routerNextV4 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: false},
			{addr: "8081", healthy: true},
			{addr: "8082", healthy: false},
			{addr: "8083", healthy: true},
			{addr: "8084", healthy: false},
			{addr: "8085", healthy: false},
		},
		curr: 3,
		next: 1,
		size: 6,
	}

	routerNextV5 = &Router{
		endpoints: []endpoint{
			{addr: "8080", healthy: true},
			{addr: "8082", healthy: false},
			{addr: "8084", healthy: false},
			{addr: "8085", healthy: true},
		},
		curr: 3,
		next: 0,
		size: 4,
	}
)
