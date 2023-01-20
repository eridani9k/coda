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

func InitializeRouter(port uint) {
	if port == 0 {
		port = 8080 // default port
	}

	// TODO: read from config file!
	balancer := NewBalancer([]string{
		"http://127.0.0.1:8081",
		"http://127.0.0.1:8082",
		"http://127.0.0.1:8083",
	})

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var forwardErr error
		var resp *http.Response
		firstRun := true // sentinel bool to ensure the following loop runs at least once.

		for forwardErr != nil || firstRun {
			firstRun = false

			endpoint, err := balancer.Advance()
			if err != nil { // ErrNoValidEndpoints
				utils.FormatMessage("No valid endpoints for routing.")
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(rw, err)
				return
			}

			targetAddr := endpoint.getAddress()

			prepareForwardingRequest(req, targetAddr)

			utils.FormatMessage(fmt.Sprintf("Routing to %s...", targetAddr))

			// FIXME: Killing endpoints often takes down others as well.
			// Things tried:
			//   - defer resp.Body.Close()
			//   - req.Close = true
			//   - using a custom http.Client()
			var httpClient = &http.Client{
				Transport: &http.Transport{},
			}

			// resp, forwardErr = http.DefaultClient.Do(req)
			resp, forwardErr = httpClient.Do(req)
			if forwardErr != nil {
				fmt.Printf("forwardErr: %s\n", err)
				endpoint.markUnhealthy()
				utils.FormatMessage(fmt.Sprintf("Error routing to %s...", targetAddr))
				continue
			}
			// defer resp.Body.Close()
		}

		// Copy all key-value pairs from the backend's
		// response Header into the new response.
		for k, values := range resp.Header {
			for _, v := range values {
				rw.Header().Set(k, v)
			}
		}

		// Copy backend's StatusCode and Body into the new response.
		rw.WriteHeader(resp.StatusCode)
		io.Copy(rw, req.Body)
	})

	utils.FormatMessage(fmt.Sprintf("Starting router on port %d...", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), proxy)
}

func prepareForwardingRequest(req *http.Request, targetAddr string) {
	url, err := url.Parse(targetAddr)
	if err != nil {
		log.Fatal(err)
	}

	req.Host = url.Host
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.RequestURI = ""

	// Set X-Forwarded-For so the backend server receives the client's address.
	originIPAddr, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Set("X-Forwarded-For", originIPAddr)

	// Supposedly resolves the issue of connections getting EOF'd when
	// they are terminated (CTRL-c) from the command line.
	// req.Close = true
}
