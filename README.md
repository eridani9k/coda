# Routing API (Round-Robin)

## Introduction

This repository contains code for an example routing API backed by a round-robin load balancing algorithm. The code was designed for the purposes of code quality review, and kept bare-bones in terms of setup and infrastructure.

_**DO NOT use this is production!**_

## Usage

A full deployment of this example consists of:
- 1 `Router` process which acts as a _reverse proxy_.
- 1..N backend `API` processes load balanced by the `Router`.

### addresses.cfg

`addresses.cfg` is a newline-delimited text file containing the list of backend addresses to be registered during the initialization of the `Router`. This file allows no other information; comments are not allowed.

The following `addresses.cfg` registers the local ports 8081, 8082, and 8083 as backend `API` processes.
```
http://127.0.0.1:8081
http://127.0.0.1:8082
http://127.0.0.1:8083
```

During `Router` initialization, each address is sent a `/ping` request to verify endpoint health. Only healthy endpoints will be added to `Router`'s load balancing algorithm.

**NOTE**: This implementation currently does not support adding backend endpoints after `Router` has completed initialization. Therefore, only the addresses in `addresses.cfg` are considered during load balancing.

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

