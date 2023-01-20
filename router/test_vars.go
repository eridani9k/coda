package router

// This file contains variables used in balancer_test.go.

var (
	emptyBalancer = &Balancer{
		endpoints: make([]*endpoint, 0),
		curr:      -1,
		size:      0,
	}

	balancerWithSingleEndpoint = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: true},
		},
		curr: 0,
		size: 1,
	}

	balancerWithMultipleEndpoints = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: true},
			{addr: ":8081", healthy: true},
			{addr: ":8082", healthy: true},
			{addr: ":8083", healthy: true},
		},
		curr: 0,
		size: 4,
	}

	balancerWithUnhealthyEndpointsV1 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: true},
			{addr: ":8081", healthy: false},
			{addr: ":8082", healthy: true},
			{addr: ":8083", healthy: false},
			{addr: ":8084", healthy: false},
			{addr: ":8085", healthy: true},
		},
		curr: 0,
		size: 6,
	}

	balancerWithUnhealthyEndpointsV2 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: false},
			{addr: ":8081", healthy: false},
			{addr: ":8082", healthy: true},
			{addr: ":8083", healthy: false},
			{addr: ":8084", healthy: true},
			{addr: ":8085", healthy: false},
		},
		curr: 0,
		size: 6,
	}

	balancerWithSingleHealthyEndpoint = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: false},
			{addr: ":8081", healthy: false},
			{addr: ":8082", healthy: false},
			{addr: ":8083", healthy: true},
			{addr: ":8084", healthy: false},
			{addr: ":8085", healthy: false},
		},
		curr: 0,
		size: 6,
	}

	balancerAllUnhealthy = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: false},
			{addr: ":8081", healthy: false},
			{addr: ":8082", healthy: false},
			{addr: ":8083", healthy: false},
			{addr: ":8084", healthy: false},
			{addr: ":8085", healthy: false},
		},
		curr: 0,
		size: 6,
	}

	balancerNextV1 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: true},
			{addr: ":8081", healthy: true},
		},
		curr: 1,
		size: 2,
	}

	balancerNextV2 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: true},
			{addr: ":8081", healthy: true},
		},
		curr: 0,
		size: 2,
	}

	balancerNextV3 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: false},
			{addr: ":8081", healthy: false},
			{addr: ":8082", healthy: false},
			{addr: ":8083", healthy: true},
			{addr: ":8084", healthy: false},
			{addr: ":8085", healthy: false},
		},
		curr: 3,
		size: 6,
	}

	balancerNextV4 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: false},
			{addr: ":8081", healthy: true},
			{addr: ":8082", healthy: false},
			{addr: ":8083", healthy: true},
			{addr: ":8084", healthy: false},
			{addr: ":8085", healthy: false},
		},
		curr: 3,
		size: 6,
	}

	balancerNextV5 = &Balancer{
		endpoints: []*endpoint{
			{addr: ":8080", healthy: true},
			{addr: ":8082", healthy: false},
			{addr: ":8084", healthy: false},
			{addr: ":8085", healthy: true},
		},
		curr: 3,
		size: 4,
	}
)
