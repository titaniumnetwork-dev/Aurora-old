package proxy

import (
	"encoding/base64"
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

var err error

// Server used for proxy
func HTTPServer(w http.ResponseWriter, r *http.Request) {
	// This will go great with json config
	blockedUserAgents := [0]string{}
	for _, userAgent := range blockedUserAgents {
		if userAgent == r.UserAgent() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, not authorized")
			return
		}
	}

	if r.TLS == nil {
		config.Scheme = "http"
	} else {
		config.Scheme = "https"
	}

	config.URL, err = url.Parse(r.URL.RequestURI())
	if err != nil {
		log.Println(err)
	}

	proxyURLB64 := config.URL.Path[len(config.HTTPPrefix):]
	proxyURLBytes, err := base64.URLEncoding.DecodeString(proxyURLB64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", err)
		log.Println(err)
		return
	}
	config.ProxyURL, err = url.Parse(string(proxyURLBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", err)
		log.Println(err)
		return
	}

	blockedDomains := [0]string{}
	for _, domain := range blockedDomains {
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

	blockedHeaders := [4]string{"Content-Security-Policy", "Content-Security-Policy-Report-Only", "Strict-Transport-Security", "X-Frame-Options"}
	for _, header := range blockedHeaders {
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
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500, %s", err)
			log.Println(err)
			return
		}
		resp.Body = respBodyInterface.(io.ReadCloser)
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
