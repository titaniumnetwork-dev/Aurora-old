package proxy

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/titaniumnetwork-dev/Aurora/modules/config"
	"github.com/titaniumnetwork-dev/Aurora/modules/rewrites"
)

// Server used for http proxy
func HTTPServer(w http.ResponseWriter, r *http.Request) {
	var err error

	for _, userAgent := range config.YAML.BlockedUserAgents {
		if userAgent == r.UserAgent() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, not authorized")
			return
		}
	}

	// TODO: Figure out how to use r.TLSs
	if config.YAML.SSLOverProxy {
		config.Scheme = "https"
	} else {
		config.Scheme = "http"
	}

	config.URL, err = url.Parse(fmt.Sprintf("%s://%s%s", config.Scheme, r.Host, r.RequestURI))
	if err != nil || config.URL.Scheme == "" || config.URL.Host == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", fmt.Sprintf("Unable to parse url, %s", config.ProxyURL.String()))
		return
	}

	proxyURLB64 := config.URL.Path[len(config.YAML.HTTPPrefix):]
	proxyURLBytes, err := base64.URLEncoding.DecodeString(proxyURLB64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", err)
		return
	}

	config.ProxyURL, err = url.Parse(string(proxyURLBytes))
	if err != nil || config.ProxyURL.Scheme == "" || config.ProxyURL.Host == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", fmt.Sprintf("Unable to parse url, %s", config.ProxyURL.String()))
		return
	}

	for _, domain := range config.YAML.BlockedDomains {
		if domain == config.ProxyURL.Hostname() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, fmt.Sprintf("401, %s has been blocked", config.ProxyURL.Hostname()))
			return
		}
	}

	for key, val := range r.Header {
		val = rewrites.Header(key, val)
		r.Header.Set(key, strings.Join(val, ", "))
	}

	tr := &http.Transport{
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", config.ProxyURL.String(), nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404, %s", err)
		return
	}

	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404, %s", err)
		return
	}
	defer resp.Body.Close()

	if config.YAML.Cap != 0 {
		http.MaxBytesReader(w, resp.Body, config.YAML.Cap)
	}

	for key, val := range resp.Header {
		val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		resp.Body = rewrites.HTML(resp.Body)
	}
	if strings.HasPrefix(contentType, "text/css") {
		respBodyInterface := rewrites.CSS(resp.Body)
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
