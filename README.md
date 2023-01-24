# Routing API (Round-Robin)

## Introduction
This repository contains code for an example routing API backed by a round-robin load balancing algorithm. The code was designed to be bare-bones in terms of setup and infrastructure for the sole purpose of code quality review.

## Usage
A full deployment of this example consists of:
- 1 `Router` process which acts as a _reverse proxy_.
- 1..N backend `API` processes load balanced by the `Router`.

### addresses.cfg
`addresses.cfg` is a newline-delimited text file containing the list of backend addresses to be registered during the initialization of the `Router`. This file allows no other information; comments are not allowed.

```
http://127.0.0.1:8081
http://127.0.0.1:8082
http://127.0.0.1:8083
```

### Running the Application
Each Golang process run can either be a `Router` or `API` server.

```golang
// From the application root directory.

// Launch a router using local port 8080.
$ go run main.go router 8080

// Launch an API backend using local port 8081.
$ go run main.go api 8081
```

There is no limit to the number of processes run, as long as local ports are available.

## Unit Tests & Coverage
```golang
$ go test -cover ./...

?   	coda		[no test files]
?   	coda/api	[no test files]
ok  	coda/router	0.096s	coverage: 28.1% of statements
ok  	coda/utils	0.081s	coverage: 50.0% of statements
```

