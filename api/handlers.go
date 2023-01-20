package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

	formatMessage(fmt.Sprintf("Starting server on port %d...", port))
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			formatMessage("Error - server closed.")
		} else {
			formatMessage(fmt.Sprintf("Error - %s", err))
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
		formatMessage("API does not support this method.")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if !json.Valid([]byte(body)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON."))
		formatMessage("Invalid JSON.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	formatMessage(fmt.Sprintf("response: %s", string(body)))
}

func home(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/home")
	fmt.Fprintf(w, "Welcome to the homepage!")
}

func formatMessage(message string) {
	fmt.Printf("[ %s ] %s\n", time.Now().Format(time.RFC3339), message)
}

func showEndpoint(endpoint string) {
	formatMessage(fmt.Sprintf("Endpoint hit: %s", endpoint))
}
