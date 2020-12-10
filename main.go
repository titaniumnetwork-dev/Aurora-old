package main

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/proxy"
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
