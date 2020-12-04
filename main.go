package main

import (
	"AuroraProxy/proxy"
	"fmt"
	"net/http"
)

const port = ":8080"

func main() {
	http.HandleFunc("/", proxy.Server)
	http.ListenAndServe(port, nil)

	fmt.Printf("Listening on port %v", port)
}
