package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func HandleRequests(port uint) {
	if port == 0 {
		port = 8080 // default port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/echo", echo)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Printf("Starting server on port %d...\n", port)
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Error - server closed.")
		} else {
			fmt.Printf("Error - %s\n", err)
		}
		os.Exit(1)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/echo")

	// API only accepts POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("API does not support this method."))
		fmt.Println("API does not support this method.")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if !json.Valid([]byte(body)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON."))
		fmt.Println("Invalid JSON.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	fmt.Printf("response: %s\\n", string(body))
}

func home(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/home")
	fmt.Fprintf(w, "Welcome to the homepage!")
}

func showEndpoint(endpoint string) {
	fmt.Printf("Endpoint hit: %s\n", endpoint)
}
