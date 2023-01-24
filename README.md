# Routing API (Round-Robin)

## Introduction
This repository contains code for an example routing API backed by a round-robin load balancing algorithm. The code was designed to be bare-bones in terms of setup and infrastructure for the sole purpose of code quality review.


## Usage
A full deployment of this example consists of:
- 1 `Router` which acts as a _reverse proxy_
- 1..N backend `API` servers load balanced by the `Router`.

Each Golang process run can either be a `Router` or `API` server.

```golang
// From the root directory

// Launch a router using local port 8080.
$ go run main.go router 8080

// Launch an API backend using local port 8081.
$ go run main.go api 8081
```

