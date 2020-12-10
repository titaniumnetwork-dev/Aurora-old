package main

import (
	"AuroraProxy/proxy"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", proxy.Server)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
