package router

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"coda/utils"
)

// Initialize() creates the reverse proxy process,
// builds an initial list of healthy backend targets,
// and serves the reverse proxy (blocking).
func Initialize(port uint, addresses []string) {
	balancer := initBalancer(addresses)

	// Set health check interval to 2 seconds.
	initHealthChecks(balancer, 2)

	// Reusable HTTP client for proxying requests.
	httpClient := &http.Client{
		Transport: &http.Transport{},
	}

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var forwardErr error
		var resp *http.Response

		// Sentinel bool to ensure the following loop runs at least once.
		firstRun := true

		for forwardErr != nil || firstRun {
			firstRun = false

			endpoint, err := balancer.Advance()
			if err != nil { // ErrNoValidEndpoints
				utils.TimestampMsg("No valid endpoints for routing.")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}

			targetAddr := endpoint.getAddress()

			prepareForwardingRequest(r, targetAddr)

			utils.TimestampMsg(fmt.Sprintf("Routing to %s...", targetAddr))

			// FIXME: Killing endpoints often takes down others as well.

			resp, forwardErr = httpClient.Do(r)
			if forwardErr != nil {
				utils.TimestampMsg(fmt.Sprintf("Error routing to %s...", targetAddr))
				endpoint.setUnhealthy()
				utils.TimestampMsg(fmt.Sprintf("Endpoint %s marked unhealthy, and will not receive requests.", endpoint.addr))
				continue
			}
			defer resp.Body.Close()
		}

		// Copy all key-value pairs from the backend's
		// response Header into the new response.
		for k, values := range resp.Header {
			for _, v := range values {
				w.Header().Set(k, v)
			}
		}

		// Copy backend's StatusCode and Body into the new response.
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, r.Body)
	})

	utils.TimestampMsg(fmt.Sprintf("Starting router on port %d...", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), proxy)
}

// initBalancer() creates an empty Balancer and registers addresses.
func initBalancer(addresses []string) *Balancer {
	balancer := NewBalancer(nil)
	for _, address := range addresses {
		register(balancer, address)
	}

	return balancer
}

// register() fires a "/ping" to the address to ensure the backend target
// is healthy. Only healthy targets are added to the Balancer.
// register() is only fired during the initialization of the Router.
func register(b *Balancer, address string) {
	err := ping(address)
	if err != nil {
		utils.TimestampMsg(fmt.Sprintf("Ping failed for %s, address will NOT be registered.", address))
		return
	}

	b.Add(address)
	utils.TimestampMsg(fmt.Sprintf("Registered address %s successfully!", address))
}

// initHealthChecks() creates an asynchronous ticker with an
// interval in seconds. Upon ticking to an interval, fires
// health checks to backend targets.
func initHealthChecks(balancer *Balancer, secondInterval int) {
	ticker := time.NewTicker(time.Duration(secondInterval) * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				healthCheck(balancer)
			}
		}
	}()
}

// healthCheck() fires a "/ping" request to all endpoints registered
// in the Balancer. This allows the Router to constantly monitor
// the health of all endpoints for optimal load balancing.
// healthCheck() logs a message if it detects a change in state,
// both from healthy -> unhealthy and vice-versa.
func healthCheck(balancer *Balancer) {
	for _, endpoint := range balancer.endpoints {
		address := endpoint.getAddress()

		err := ping(address)
		if err != nil {
			utils.TimestampMsg(fmt.Sprintf("Health check failed for %s.", address))
			endpoint.setUnhealthy()
			continue
		}

		// Notify if endpoint is recovering from an unhealthy state.
		if !endpoint.isHealthy() {
			utils.TimestampMsg(fmt.Sprintf("%s recovering from an unhealthy state.", address))
			balancer.ResetCurr()
		}

		endpoint.setHealthy()
	}
}

// ping() fires a "/ping" request to the address.
// The response is ignored since the backend returns HTTP 200 by default.
func ping(address string) error {
	_, err := http.Get(fmt.Sprintf("%s/ping", address))
	return err
}

// prepareForwardingRequest() primes the request by modifying
// the necessary headers required for proxying to the backend.
func prepareForwardingRequest(r *http.Request, targetAddr string) {
	url, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatal(err)
	}

	r.Host = url.Host
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.RequestURI = ""

	// Set X-Forwarded-For so the backend server receives the origin's address.
	originIPAddr, _, _ := net.SplitHostPort(r.RemoteAddr)
	r.Header.Set("X-Forwarded-For", originIPAddr)

	// Attempt to resolve the issue of connections getting EOF'd when
	// they are terminated (CTRL-c) from the command line.
	r.Close = true
}
