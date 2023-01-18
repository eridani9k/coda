package router

import (
	//"io"
	//"net/http"
	//"net/http/httputil"
)

func InitializeRouter(port uint) {
	if port == 0 {
		port = 8080 // default port
	}

	//mux := http.NewServeMux()
	//mux.HandleFunc("/postjson", postJSON)
}
