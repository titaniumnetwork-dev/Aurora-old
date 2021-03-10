package main

import (
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"github.com/titaniumnetwork-dev/Aurora/modules/proxy"
	"log"
	"net/http"
	"os"
)

var err error

func main() {
	config.HTTPPrefix, config.HTTPPrefixExists = os.LookupEnv("HTTPPREFIX")
	config.WSPrefix, config.WSPrefixExists = os.LookupEnv("WSPREFIX")
	// config.WRTCPrefix, config.WRTCPrefixExists = os.LookupEnv("WRTCPREFIX")
	if config.HTTPPrefixExists == false {
		log.Fatal("You need to specify an http prefix")
	} else if config.WSPrefixExists == false {
		log.Fatal("You need to specify an ws prefix")
	} else {
		http.HandleFunc(config.HTTPPrefix, proxy.HTTPServer)
		http.HandleFunc(config.WSPrefix, proxy.WSServer)
		// http.HandleFunc(config.WRTCPrefix, proxy.WRTCServer)
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	config.Port, config.PortExists = os.LookupEnv("PORT")
	if config.PortExists {
		config.SSLCert, config.SSLCertExists = os.LookupEnv("CERT")
		config.SSLKey, config.SSLKeyExists = os.LookupEnv("KEY")

		if config.SSLCertExists && config.SSLKeyExists {
			err = http.ListenAndServeTLS(config.Port, config.SSLCert, config.SSLKey, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = http.ListenAndServe(config.Port, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("You need to specify a port")
	}
}
