package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"coda/utils"
)

// Initialize() creates the API server multiplexer
// and registers several URL endpoint patterns to their
// respective handling functions.
func Initialize(port uint) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/ping", ping)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	utils.TimestampMsg(fmt.Sprintf("Starting server on port %d...", port))
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			utils.TimestampMsg("Error - server closed.")
		} else {
			utils.TimestampMsg(fmt.Sprintf("Error - %s", err))
		}
		os.Exit(1)
	}
}

// echo() responds with the same body as the request.
// The request must have the POST method and must contain
// a valid JSON body.
func echo(w http.ResponseWriter, r *http.Request) {
	// TODO: Find a way to utilize r.Context() instead of using showEndpoint().
	// showEndpoint("/echo")

	// API only accepts POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("API does not support this method."))
		utils.TimestampMsg("API does not support this method.")
		return
	}

	// FIXME: EOF propagation issue with request body being closed
	// too early when OTHER servers are killed via CTRL-c.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unexpected EOF."))
		utils.TimestampMsg("Unexpected EOF.")
		return
	}

	if !json.Valid([]byte(body)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON."))
		utils.TimestampMsg("Invalid JSON.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	utils.TimestampMsg(fmt.Sprintf("response: %s", string(body)))
}

// ping() allows heartbeart checking and returns 200.
func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func home(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/home")
	fmt.Fprintf(w, "Welcome to the homepage!")
}

// showEndpoint() prints a string to STDOUT.
// Used as a debug utility.
func showEndpoint(endpoint string) {
	utils.TimestampMsg(fmt.Sprintf("Endpoint hit: %s", endpoint))
}
