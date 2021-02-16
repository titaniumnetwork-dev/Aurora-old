package main

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	// TODO: Add a proxy main page at root of path
	// TODO: Figure out how to pass config to a container
	// Path doesn't work?
	path, pathExists := os.LookupEnv("PROXYPATH")
	if pathExists {
		http.HandleFunc(path, proxy.Server)
	} else {
		log.Fatal("You need to specify a path")
	}

	port, portExists := os.LookupEnv("PROXYPORT")
	if portExists {
		sslCert, sslCertExists := os.LookupEnv("CERTPATH")
		sslKey, sslKeyExists := os.LookupEnv("KEYPATH")
		if sslCertExists && sslKeyExists {
			err := http.ListenAndServeTLS(port, sslCert, sslKey, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := http.ListenAndServe(port, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("You need to specify a port")
	}
}
