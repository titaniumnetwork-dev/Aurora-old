package main

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/global"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/proxy"
	"log"
	"net/http"
	"os"
)

func main() {
	// TODO: Add a proxy main page at root of path
	// TODO: Figure out how to pass config to a container
	// TODO: Change env name to prefix
	global.Prefix, global.PrefixExists = os.LookupEnv("PREFIX")
	if global.PrefixExists {
		http.HandleFunc(global.Prefix, proxy.Server)
		http.Handle("/", http.FileServer(http.Dir("./static")))
	} else {
		log.Fatal("You need to specify a prefix")
	}

	// TODO: Change env name to PORT
	global.Port, global.PortExists = os.LookupEnv("PORT")
	if global.PortExists {
		if global.SSLCertExists && global.SSLKeyExists {
			global.SSLCert, global.SSLCertExists = os.LookupEnv("CERTPATH")
			global.SSLKey, global.SSLKeyExists = os.LookupEnv("KEYPATH")
			err := http.ListenAndServeTLS(global.Port, global.SSLCert, global.SSLKey, nil)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := http.ListenAndServe(global.Port, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("You need to specify a port")
	}
}
