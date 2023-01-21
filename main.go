package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"coda/api"
	"coda/router"
)

const (
	serverTypeRouter = "router"
	serverTypeAPI    = "api"
)

var (
	// Retrieves application root directory.
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
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
		endpoints := readEndpointFile("endpoints.cfg")
		router.InitializeRouter(uint(port), endpoints)
	} else {
		api.HandleRequests(uint(port))
	}
}

func readEndpointFile(configFile string) []string {
	file, err := os.Open(filepath.Join(basepath, configFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var endpoints []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		endpoints = append(endpoints, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return endpoints
}
