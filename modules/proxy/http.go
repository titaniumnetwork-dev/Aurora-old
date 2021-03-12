package proxy

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"github.com/titaniumnetwork-dev/Aurora/modules/rewrites"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Server used for http proxy
func HTTPServer(w http.ResponseWriter, r *http.Request) {
	var err error

	for _, userAgent := range config.BlockedUserAgents {
		if userAgent == r.UserAgent() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, not authorized")
			return
		}
	}

	config.SSLOverProxy, config.SSLOverProxyExists = os.LookupEnv("SSLOVERPROXY")
	if config.SSLOverProxyExists == false {
		config.SSLOverProxy == false
	}
	if r.TLS != nil || config.SSLOverProxy == true {
		config.HTTPScheme = "https"
	} else {
		config.HTTPScheme = "http"
	}

	config.URL, err = url.Parse(config.Scheme + "://" + r.Host + r.RequestURI)
	if err != nil || config.URL.Scheme == "" || config.URL.Host == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", errors.New("Unable to parse url"))
		return
	}

	proxyURLB64 := config.URL.Path[len(config.HTTPPrefix):]
	proxyURLBytes, err := base64.URLEncoding.DecodeString(proxyURLB64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", err)
		return
	}

	config.ProxyURL, err = url.Parse(string(proxyURLBytes))
	if err != nil || config.ProxyURL.Scheme == "" || config.ProxyURL.Host == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", errors.New("Unable to parse url"))
		return
	}

	for _, domain := range config.BlockedDomains {
		if domain == config.ProxyURL.Hostname() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, this domain has been blocked")
			return
		}
	}

	// TODO: Add the option to cap file transfer size with environment variable
	tr := &http.Transport{
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", config.ProxyURL.String(), nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404, %s", err)
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404, %s", err)
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	for _, header := range config.BlockedHeaders {
		delete(resp.Header, header)
	}
	for key, val := range resp.Header {
		val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		resp.Body, err = rewrites.HTML(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500, %s", err)
			log.Println(err)
			return
		}
	}
	if strings.HasPrefix(contentType, "text/css") {
		respBodyInterface, err := rewrites.CSS(resp.Body)
		resp.Body = respBodyInterface.(io.ReadCloser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500, %s", err)
			log.Println(err)
			return
		}
	}
	// Currently low priority
	/*
		if strings.HasPrefix(contentType, "image/svg") {
			resp.Body = rewrites.SVG(resp.Body)
		}
		if strings.HasPrefix(contentType, "application/xml") strings.HasPrefix(contentType, "text/xml") {
			resp.Body = rewrites.XML(resp.Body)
		}
	*/

	io.Copy(w, resp.Body)
}
