# Routing API (Round-Robin)

## Introduction

This repository contains code for an example routing API backed by a round-robin load balancing algorithm. The code was designed for the purposes of code quality review, and kept bare-bones in terms of setup and infrastructure.

_**This is sample code. DO NOT deploy this in a production environment!**_

## Usage

### addresses.cfg

`addresses.cfg` is a newline-delimited text file containing the list of backend addresses to be registered during the initialization of the `Router`. This file allows no other information; comments are not allowed.

The following `addresses.cfg` registers the local ports 8081, 8082, and 8083 as backend `API` processes.
```
http://127.0.0.1:8081
http://127.0.0.1:8082
http://127.0.0.1:8083
```

During `Router` initialization, each address is sent a `/ping` request to verify endpoint health. Only healthy endpoints will be added to the `Router`'s load balancing algorithm.

**NOTE**: This implementation currently does not support adding backend endpoints after `Router` has completed initialization. Therefore, only the addresses in `addresses.cfg` are considered during load balancing.


### Running the Application

Each Golang process can be run either as a `Router` or an `API` server.

```golang
// Executed from the application root directory.

// Launch an API backend using local port 8081.
$ go run main.go api 8081

// Launch a router using local port 8080.
$ go run main.go router 8080
```

There is no limit to the number of processes run, as long as local ports are available.


### Deployment Example

This example consists of the following deployments and showcases basic round-robin routing:
- 1 `Router` process which acts as a _reverse proxy_.
- 1..N backend `API` processes load balanced by the `Router`.

In this example: 
- 1 `Router` process is launched on local port 8080.
- 3 `API` processes are launched on local ports 8081, 8082, and 8083. These 3 addresses were defined in `addresses.cfg`.

The `API` processes should be launched before the `Router` since all endpoints are pinged for health before registration into the load balancing algorithm.

```golang
// Executed from the application root directory.
// Each command in this block is launched in a separate
// terminal window as all processes are blocking.

// Launching the API processes.
$ go run main.go api 8081
$ go run main.go api 8082
$ go run main.go api 8083

/* Expected output:
[ 2023-01-24T18:52:07+08:00 ] Starting server on port 8081...
*/

// Launching the Router process.
$ go run main.go router 8080

/* Expected output:
[ 2023-01-24T18:53:44+08:00 ] Registered address http://127.0.0.1:8081 successfully!
[ 2023-01-24T18:53:44+08:00 ] Registered address http://127.0.0.1:8082 successfully!
[ 2023-01-24T18:53:44+08:00 ] Registered address http://127.0.0.1:8083 successfully!
[ 2023-01-24T18:53:44+08:00 ] Starting router on port 8080...
*/
```

Run `curl.sh <ROUTER_PORT>` to start sending requests to the `Router` process.
`curl.sh` runs 100 iterations of `curl POST ...` to simulate incoming requests.

```bash
$ ./curl.sh 8080
```

The `Router` should now print routing destinations in a round-robin basis:
```
[ 2023-01-24T19:02:48+08:00 ] Routing to http://127.0.0.1:8081...
[ 2023-01-24T19:02:51+08:00 ] Routing to http://127.0.0.1:8082...
[ 2023-01-24T19:02:54+08:00 ] Routing to http://127.0.0.1:8083...
[ 2023-01-24T19:02:57+08:00 ] Routing to http://127.0.0.1:8081...
[ 2023-01-24T19:03:00+08:00 ] Routing to http://127.0.0.1:8082...
[ 2023-01-24T19:03:03+08:00 ] Routing to http://127.0.0.1:8083...

...
```

## Unit Tests & Coverage

```golang
$ go test -cover ./...

?   	coda		[no test files]
?   	coda/api	[no test files]
ok  	coda/router	0.096s	coverage: 28.1% of statements
ok  	coda/utils	0.081s	coverage: 50.0% of statements
```

