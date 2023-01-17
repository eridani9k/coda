package main

import (
	"fmt"
	"os"
	"strconv"

	"coda/api"
)

func main() {
	var port uint
	if len(os.Args) == 2 {
		p, err := strconv.ParseUint(os.Args[1], 10, 0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		port = uint(p)
	} else {
		port = 8080
	}

	api.HandleRequests(port)
}
