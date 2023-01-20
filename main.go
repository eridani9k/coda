package main

import (
	"fmt"
	"os"
	"strconv"

	"coda/api"
	"coda/router"
)

const (
	serverTypeRouter = "router"
	serverTypeAPI    = "api"
)

func main() {
	// Parse command line arguments.
	//   os.Args[1] - Server type. Valid options: router | api
	//   os.Args[2] - Local port. Valid options: 1-65535

	if len(os.Args) != 3 {
		fmt.Println("Incorrect number of arguments, terminating.")
		os.Exit(1)
	}

	serverType := os.Args[1]
	if serverType != serverTypeRouter && serverType != serverTypeAPI {
		fmt.Println("Invalid server type. Valid options: 'router' | 'api'")
		os.Exit(1)
	}

	port, err := strconv.ParseUint(os.Args[2], 10, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if port < 1 || port > 65535 {
		fmt.Println("Invalid port. Valid options: 1-65535")
		os.Exit(1)
	}

	if serverType == serverTypeRouter {
		router.InitializeRouter(uint(port))
	} else {
		api.HandleRequests(uint(port))
	}
}
