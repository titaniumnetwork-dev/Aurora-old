package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"

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
		// There is a security bug here where the other users sites they are accessing will get leaked
		http.HandleFunc(config.YAML.HTTPPrefix, proxy.HTTPServer)
		http.HandleFunc(config.YAML.WSPrefix, proxy.WSServer)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			type d struct {
				HTTPPrefix string
				WSPrefix   string
				URL        *url.URL
				ProxyURL   *url.URL
			}

			data := d{}

			data.HTTPPrefix = config.YAML.HTTPPrefix
			data.WSPrefix = config.YAML.WSPrefix
			data.URL = config.URL
			data.ProxyURL = config.ProxyURL

			switch r.URL.Path {
			case "/":
				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-Type", "text/html")
				tmpl, err := template.ParseFiles("static/index.html")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, fmt.Sprintf("500, %s", err))
					return
				}
				err = tmpl.Execute(w, data)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, fmt.Sprintf("500, %s", err))
					return
				}
			case "/inject":
				w.Header().Add("Content-Type", "application/javascript")
				w.WriteHeader(http.StatusOK)
				tmpl, err := template.ParseFiles("static/inject.js")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, fmt.Sprintf("500, %s", err))
					return
				}
				err = tmpl.Execute(w, data)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, fmt.Sprintf("500, %s", err))
					return
				}
			case "/robots.txt":
				// TODO: Handle this
			default:
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, fmt.Sprintf("401, can't get %s", r.URL.Path))
			}
		})
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
