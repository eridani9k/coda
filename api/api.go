package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleRequests(port uint) {
	if port == 0 {
		port = 8080 // default port
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)

	fmt.Printf("Starting server on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/echo")

	// API only accepts POST
	if r.Method != http.MethodPost {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// if JSON is invalid return an error code.
	if v := json.Valid([]byte(body)); !v {
		fmt.Println("invalid JSON!")
		return
	}

	fmt.Println(string(body))
}

func home(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/home")
	fmt.Fprintf(w, "Welcome to the homepage!")
}

func showEndpoint(endpoint string) {
	fmt.Printf("Endpoint hit: %s\n", endpoint)
}
