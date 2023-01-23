package router

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"coda/utils"
)

func InitializeRouter(port uint, addresses []string) {
	balancer := NewBalancer(nil)
	for _, address := range addresses {
		register(balancer, address)
	}

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var forwardErr error
		var resp *http.Response
		firstRun := true // sentinel bool to ensure the following loop runs at least once.

		for forwardErr != nil || firstRun {
			firstRun = false

			endpoint, err := balancer.Advance()
			if err != nil { // ErrNoValidEndpoints
				utils.FormatMessage("No valid endpoints for routing.")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}

			targetAddr := endpoint.getAddress()

			prepareForwardingRequest(r, targetAddr)

			utils.FormatMessage(fmt.Sprintf("Routing to %s...", targetAddr))

			// FIXME: Killing endpoints often takes down others as well.
			// Things tried:
			//   - defer resp.Body.Close()
			//   - req.Close = true
			//   - using a custom http.Client()
			httpClient := &http.Client{
				Transport: &http.Transport{},
			}

			resp, forwardErr = httpClient.Do(r)
			if forwardErr != nil {
				fmt.Println("forwardErr: ", err)
				utils.FormatMessage(fmt.Sprintf("Error routing to %s...", targetAddr))
				endpoint.markUnhealthy()
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

	utils.FormatMessage(fmt.Sprintf("Starting router on port %d...", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), proxy)
}

func register(b *Balancer, address string) {
	_, err := http.Get(fmt.Sprintf("%s/ping", address))
	if err != nil {
		utils.FormatMessage(fmt.Sprintf("Ping failed for %s, address will NOT be registered.", address))
		return
	}

	b.Add(address)
	utils.FormatMessage(fmt.Sprintf("Registered address %s successfully!", address))
}

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

	// Supposedly resolves the issue of connections getting EOF'd when
	// they are terminated (CTRL-c) from the command line.
	r.Close = true
}
