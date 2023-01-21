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

func HandleRequests(port uint) {
	if port == 0 {
		port = 8080 // default port
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/ping", ping)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	utils.FormatMessage(fmt.Sprintf("Starting server on port %d...", port))
	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			utils.FormatMessage("Error - server closed.")
		} else {
			utils.FormatMessage(fmt.Sprintf("Error - %s", err))
		}
		os.Exit(1)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/echo")

	// TODO: Find a way to utilize r.Context() instead of using
	// showEndpoint().
	//fmt.Println(r.Context())

	// API only accepts POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("API does not support this method."))
		utils.FormatMessage("API does not support this method.")
		return
	}

	// FIXME: EOF issue with request body being closed too early.
	// Might be fixed after markUnhealthy: Post "http://127.0.0.1:8082/echo": readfrom tcp 127.0.0.1:52729->127.0.0.1:8082: http: invalid Read on closed Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("io.ReadAll error: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unexpected EOF."))
		utils.FormatMessage("Unexpected EOF.")
		return
	}

	if !json.Valid([]byte(body)) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON."))
		utils.FormatMessage("Invalid JSON.")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	utils.FormatMessage(fmt.Sprintf("response: %s", string(body)))
}

// ping allows heartbeart checking.
func ping(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/ping")

	w.WriteHeader(http.StatusOK)
}

func home(w http.ResponseWriter, r *http.Request) {
	showEndpoint("/home")
	fmt.Fprintf(w, "Welcome to the homepage!")
}

func showEndpoint(endpoint string) {
	utils.FormatMessage(fmt.Sprintf("Endpoint hit: %s", endpoint))
}
