package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
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

	addressFile = "addresses.cfg"
)

var (
	// Retrieves application root directory.
	// Used to resolve address configuration file.
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// main() parses and handles two types of command line arguments.
//
//	os.Args[1] - Server type. Valid options: router | api
//	os.Args[2] - Local port. Valid options: 1-65535
//
// Examples:
//  1. To start a router process at port 8080:
//     $ go run main.go router 8080
//  2. To start an API backend process at port 8081:
//     $ go run main.go api 8081
//
// In general, a deployment set consists of a single Router
// fronting 1..N backend processes.
func main() {
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
		router.InitRouter(uint(port), readAddressFile(addressFile))
	} else {
		api.HandleRequests(uint(port))
	}
}

// readAddressFile() reads a line-separated text file containing
// the addresses of each backend endpoint. Each address needs to
// be fully qualified (e.g. scheme, domain, port).
// The text file being read should not contain any other types of
// information aside from addresses. Comments are not allowed.
func readAddressFile(configFile string) []string {
	file, err := os.Open(filepath.Join(basepath, configFile))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var addresses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()

		_, err := url.Parse(t)
		if err != nil {
			fmt.Printf("Malformed URL: %s. Ignoring entry.\n", t)
			continue
		}

		addresses = append(addresses, t)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return addresses
}
