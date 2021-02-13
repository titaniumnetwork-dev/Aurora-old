package main

import (
	//	"github.com/titaniumnetwork-dev/AuroraProxy/modules/config"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	// TODO: Add a proxy main page at root of path
	// TODO: Add optional support ssl enabled through environment variable
	path, pathExists := os.LookupEnv("PROXYPATH")
	if pathExists {
		http.HandleFunc(path, proxy.Server)
	} else {
		http.HandleFunc("/", proxy.Server)
	}

	port, portExists := os.LookupEnv("PROXYPORT")
	if portExists {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// TODO: Once optional ssl support is added change to port 443 instead as default
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
