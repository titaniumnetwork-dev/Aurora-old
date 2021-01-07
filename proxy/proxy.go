package proxy

import (
	// "github.com/titaniumnetwork-dev/AuroraProxy/rewrites"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	re := regexp.MustCompile(`(\:\/)([^\/])`)
	url := re.ReplaceAllString(r.URL.Path[1:], "$1/$2")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	for key, val := range resp.Header {
		w.Header().Set(key, strings.Join(val, ", "))
	}
	// TODO: Add header rewriting via regex
	w.WriteHeader(resp.StatusCode)

	// Theoretical code in preperation for rewrites
	/*
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") { // Uses html parsing with the (experimental library) html
			rewrittenBody := Rewrites.Html(resp.Body)
		}
		if else strings.HasPrefix(contentType, "text/css") { // Uses regular expressions with the library regexp
			rewrittenBody := Rewrites.Css(resp.Body)
		}
		if else strings.HasPrefix(contentType, "text/js") { // Uses regular expressions with the library regexp
			rewrittenBody := Rewrites.Js(resp.Body)
			// TODO: Add js injection code here
		}

		io.Copy(w, rewrittenBody)
	*/

	io.Copy(w, resp.Body)
}
