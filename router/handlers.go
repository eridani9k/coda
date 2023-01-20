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

const (
	TESTAPI = "http://127.0.0.1:8081"
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

	// url should be the *URL object of the backend target.
	url, err := url.Parse(TESTAPI)
	if err != nil {
		log.Fatal(err)
	}

	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("Populating proxy request with backend destination...")
		req.Host = url.Host
		req.URL.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.RequestURI = ""

		// Set X-Forwarded-For so the backend server receives the client's address.
		// TODO: verify this
		fmt.Println("setting X-Forwarded-For header...")
		// This returns server IP. Port is ignored.
		s, _, _ := net.SplitHostPort(req.RemoteAddr)
		fmt.Println(s)
		req.Header.Set("X-Forwarded-For", s)

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
