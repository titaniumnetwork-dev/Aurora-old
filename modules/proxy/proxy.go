package proxy

import (
	//	"os"
	"encoding/base64"
	"fmt"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/rewrites"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Server used for proxy
// TODO: Add an optional header blocklist to block site blockers and crawlers
// TODO: Use seed based url obfustication
func Server(w http.ResponseWriter, r *http.Request) {
	// TODO: Add the option to cap file transfer size with environment variable
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	proxyUriBytes, err := base64.StdEncoding.DecodeString(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: Send get error page w/ error template page (get path from environment variable)
		fmt.Fprintf(w, "Placeholder error")
		log.Println(err)
		return
	}
	proxyUri := string(proxyUriBytes)
	proxyUri = rewrites.ProxyUri(proxyUri)

	req, err := http.NewRequest("GET", proxyUri, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: Send get error page w/ error template page (get path from environment variable)
		fmt.Fprintf(w, "Placeholder error")
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO: Send get error page w/ error template page (get path from environment variable)
		fmt.Fprintf(w, "Placeholder error")
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// TODO: Remove CORS blocking headers
	for key, val := range resp.Header {
		val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	// TODO: Add more content type checking due to there being alternatives used on the web
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		resp.Body = rewrites.Html(resp.Body)
	}
	/*
		if strings.HasPrefix(contentType, "text/css") {
			resp.Body = rewrites.Css(resp.Body)
		}
		if strings.HasPrefix(contentType, "text/javascript") {
			resp.Body = rewrites.Js(resp.Body)
		}
	*/
	// Currently low priority
	/*
		if strings.HasPrefix(contentType, "text/xml") {
			resp.Body = rewrites.Xml(resp.Body)
		}
	*/

	io.Copy(w, resp.Body)
}
