package proxy

import (
	"encoding/base64"
	"fmt"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/global"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/rewrites"
	"io"
	"log"
	"net/http"
	//	"net/url"
	"strings"
	"time"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	// TODO: Add header blacklist
	// blockedHeaders := [0]string{}
	
	global.Host = r.Host

	proxyURIB64 := r.URL.Path[len(global.Prefix):]
	proxyURIBytes, err := base64.URLEncoding.DecodeString(proxyURIB64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		log.Println(err)
		return
	}
	global.ProxyURI = string(proxyURIBytes)

	/*
	proxyURIParsed, err := url.Parse(global.ProxyURI)
	if err != nil {
		log.Println(err)
	}
	*/

	// TODO: Add the option to cap file transfer size with environment variable
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", global.ProxyURI, nil)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err)
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	blockedHeaders := [2]string{"Content-Security-Policy", "Strict-Transport-Security"}
	for i := 0; i < len(blockedHeaders); i++ {
		delete(resp.Header, blockedHeaders[i])
	}
	for key, val := range resp.Header {
		val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	// TODO: Add more content type checking due to there being alternatives used on the web
	// TODO: Not being checked correctly
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
