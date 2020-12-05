package main

import (
	"AuroraProxy/proxy"
	"fmt"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	fmt.Printf("Starting server on port %v\n", port[1:])

	http.HandleFunc("/", proxy.Server)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal(err)
	}
}
