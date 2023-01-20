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

	// ReverseProxy
	//   1. Create a reverse proxy handler.
	//   2. set backend server/s for retrieval
	//   3. rewrite appropriate headers.
	//   4. send request to backend server/s
	//   5. edit headers before sending response back to client?!`

	// TODO: read from config file!
	balancer := NewBalancer([]string{
		"http://127.0.0.1:8081",
		"http://127.0.0.1:8082",
	})

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		endpoint, err := balancer.Advance()
		if err != nil {
			log.Fatal(err)
		}

		backendAddr := endpoint.getAddress()
		url, err := url.Parse(backendAddr)
		if err != nil {
			log.Fatal(err)
		}

		req.Host = url.Host
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.RequestURI = ""

		// Set X-Forwarded-For so the backend server receives the client's address.
		// TODO: verify this
		// This returns server IP. Port is ignored.
		originIPAddr, _, _ := net.SplitHostPort(req.RemoteAddr)
		req.Header.Set("X-Forwarded-For", originIPAddr)

		utils.FormatMessage(fmt.Sprintf("Routing to %s...", backendAddr))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
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
