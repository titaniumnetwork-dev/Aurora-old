package proxy

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/global"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/rewrites"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	blockedUserAgents := [0]string{}
	for _, userAgent := range blockedUserAgents {
		if userAgent == r.UserAgent() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401 not allowed")
			return
		}
	}

	if r.TLS != nil {
		global.Scheme = "https://"
	} else {
		global.Scheme = "http://"
	}

	global.Host = r.Host

	proxyURIB64 := r.URL.Path[len(global.Prefix):]
	proxyURIBytes, err := base64.URLEncoding.DecodeString(proxyURIB64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", err)
		log.Println(err)
		return
	}
	global.ProxyURI = string(proxyURIBytes)

	// TODO: Add the option to cap file transfer size with environment variable
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", global.ProxyURI, nil)
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

	// TODO: Block more headers
	blockedHeaders := [2]string{"Content-Security-Policy", "Strict-Transport-Security"}
	for i := 0; i < len(blockedHeaders); i++ {
		delete(resp.Header, blockedHeaders[i])
	}
	for key, val := range resp.Header {
		val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	// TODO: Rewrite audio/video metadata for streams
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		resp.Body = rewrites.HTML(resp.Body)
	}
	if strings.HasPrefix(contentType, "text/css") {
		resp.Body = rewrites.CSS(resp.Body)
	}
	if strings.HasPrefix(contentType, "application/javascript") {
		resp.Body = rewrites.JS(resp.Body)
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
