package api

import (
	"fmt"
	"log"
	"net/http"
)

func HandleRequests(port uint) {
	if port == 0 {
		port = 8080 // default port
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the homepage!")
	showEndpoint("/home")
}

func echo(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/echo")
	return
}

func showEndpoint(endpoint string) {
	fmt.Printf("Endpoint hit: %s\n", endpoint)
}
