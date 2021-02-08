package proxy

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/rewrites"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Server used for proxy
// TODO: Add error catching so the program doesn't crash immediately
func Server(w http.ResponseWriter, r *http.Request) {
	// TODO: Add the option to cap file transfer size once configuration is supported
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	// TODO: Add the option to block user agents and send them to a blocked page once configuration is supported
	url := rewrites.ProxyUrl(r.URL.Path[1:])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// TODO: Remove CORS headers
	for key, val := range resp.Header {
		// val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	// TODO: Add more content type checking due to there being alternatives used on the web
	/*
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			body = rewrites.Html(resp.Body)
		}
		if strings.HasPrefix(contentType, "text/css") {
			body = rewrites.Css(resp.Body)
		}
		if strings.HasPrefix(contentType, "text/javascript") {
			// body = rewrites.Js(resp.Body)
			body = rewrites.JsInject(resp.Body)
		}
	*/
	/*
		if strings.HasPrefix(contentType, "text/xml") {
			body = rewrites.Xml(resp.Body)
		}
	*/

	io.Copy(w, resp.Body)
}
