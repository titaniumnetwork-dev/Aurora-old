package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"github.com/titaniumnetwork-dev/Aurora/modules/proxy"
	"gopkg.in/yaml.v2"
)

var err error

func main() {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("You need to add a config")
	}

	err = yaml.Unmarshal(file, &config.YAML)

	if config.YAML.HTTPPrefix == "" {
		log.Fatal("You need to specify an http prefix")
	} else if config.YAML.WSPrefix == "" {
		log.Fatal("You need to specify an ws prefix")
	} else {
		http.HandleFunc(config.YAML.HTTPPrefix, proxy.HTTPServer)
		http.HandleFunc(config.YAML.WSPrefix, proxy.WSServer)
		// TODO: Serve main.html and use templates
		http.Handle("/", http.FileServer(http.Dir("./static")))
	}

	if config.YAML.Port == "" {
		log.Fatal("You need to specify a port")
	}

	if config.YAML.Cert != "" && config.YAML.Key != "" {
		err = http.ListenAndServeTLS(config.YAML.Port, config.YAML.Cert, config.YAML.Key, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = http.ListenAndServe(config.YAML.Port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
