package main

import (
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"github.com/titaniumnetwork-dev/Aurora/modules/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	config.HTTPPrefix, config.HTTPPrefixExists = os.LookupEnv("HTTPPREFIX")
	config.WSPrefix, config.WSPrefixExists = os.LookupEnv("WSPREFIX")
	if !config.HTTPPrefixExists {
		log.Fatal("You need to specify an http prefix")
	} else if !config WSPrefixExists {
		log.Fatal("You need to specify an ws prefix")
	} else {
		http.HandleFunc(config.HTTPPrefix, http.Server)
		http.HandleFunc(config.WSPrefix, ws.Server)
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	config.Port, config.PortExists = os.LookupEnv("PORT")
	if config.PortExists {
		if config.SSLCertExists && config.SSLKeyExists {
			config.SSLCert, config.SSLCertExists = os.LookupEnv("CERT")
			config.SSLKey, config.SSLKeyExists = os.LookupEnv("KEY")

			err := http.ListenAndServeTLS(config.Port, config.SSLCert, config.SSLKey, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := http.ListenAndServe(config.Port, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("You need to specify a port")
	}
}
