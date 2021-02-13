package proxy

import (
	//	"os"
	"fmt"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/rewrites"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	// TODO: See if I can do this in main.go instead
	/*
		homePagePath, homePageExists := os.LookupEnv("HOMEPAGEPATH")
		homePagePath = "././" + homePagePath
		if r.URL.Path[1:] == "/" || homePageExists {
			// TODO: Read proxy home page file
			io.Copy(w, homePagePath)
			w.WriteHeader(200)
			return
		}

		// Figure out how to get these variables to rewrites.go maybe make an environment variable
	 	Domain := r.URL
		// TODO: Make a uri variable
		URI :=
	*/

	// TODO: Add the option to cap file transfer size with environment variable
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	// TODO: Add an optional header blocklist to block site blockers and crawlers (get info from blocklist.json file)
	url := rewrites.ProxyUrl(r.URL.Path[1:])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// TODO: Send get error page w/ error template page (get path from environment variable)
		fmt.Fprintf(w, "Placeholder error")
		// TODO: Add status code header and close the response writer
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		// TODO: Send get error page w/ error template page (get path from environment variable)
		fmt.Fprintf(w, "Placeholder error")
		// TODO: Add status code header and close the response write
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
	/*
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			resp.Body = rewrites.Html(resp.Body)
		}
		if strings.HasPrefix(contentType, "text/css") {
			body = rewrites.Css(resp.Body)
		}
		if strings.HasPrefix(contentType, "text/javascript") {
			body = rewrites.Js(resp.Body)
		}
	*/
	// Currently low priority
	/*
		if strings.HasPrefix(contentType, "text/xml") {
			body = rewrites.Xml(resp.Body)
		}
	*/

	io.Copy(w, resp.Body)
}
