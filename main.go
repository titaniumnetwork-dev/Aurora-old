package main

import (
	//	"github.com/titaniumnetwork-dev/AuroraProxy/config"
	"github.com/titaniumnetwork-dev/AuroraProxy/proxy"
	"log"
	"net/http"
)

func main() {
	// TODO: Add custom path once configuration is supported
	http.HandleFunc("/", proxy.Server)

	// TODO: Add custom port once configuration is supported
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
