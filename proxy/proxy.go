package proxy

import (
	//	"github.com/titaniumnetwork-dev/AuroraProxy/rewrites"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Server used for proxy
// TODO: If the user agent is a site blocker send a 404
func Server(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

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
		val = rewrites.Header(key, val)
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		body = Rewrites.Html(resp.Body)
	}
	if strings.HasPrefix(contentType, "text/js") {
		body = Rewrites.Js(resp.Body)
	}

	io.Copy(w, resp.Body)
}
