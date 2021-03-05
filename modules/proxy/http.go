package proxy

import (
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/config"
	"github.com/titaniumnetwork-dev/AuroraProxy/modules/rewrites"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
	"strings"
	"time"
	"net/url"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	// This will go great with json config
	blockedUserAgents := [0]string{}
	for _, userAgent := range blockedUserAgents {
		if userAgent == r.UserAgent() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, not authorized")
			return
		}
	}

	global.Cookie, global.CookieExists := os.LookupEnv("COOKIE")
	global.Cookie = strings.Split(global.Cookie, "=")
	if global.CookieExists && len(global.Cookie) == 2 {
		cookie, err := http.Cookie(global.Cookie[0])
		// Yeah this can't be cookie.name it has to be something different for value
		if err != nil || cookie.name != global.Cookie[1] {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, not authorized")
			return
		}
	}

	if r.TLS == nil {
		global.Scheme = "http:"
	} else {
		global.Scheme = "https:"
	}

	global.URL, err = url.Parse(req.URL.RequestURI())
	if err != nil {
		log.Println(err)
	}

	proxyURLB64 := global.URL.Path[len(global.Prefix):]
	proxyURLBytes, err := base64.URLEncoding.DecodeString(proxyURLB64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500, %s", err)
		log.Println(err)
		return
	}
	global.ProxyURL = url.Parse(string(proxyURLBytes))

	// This will go great with json config
	blockedDomains := [0]string{}
	for _, domain := range blockedDomains {
		if domain == global.ProxyURL.Hostname() {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "401, this domain has been blocked")
			return
		}
	}

	// TODO: Add the option to cap file transfer size with environment variable
	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", global.ProxyURL.String(), nil)
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

	// This will go great with json config
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
		resp.Body, err = rewrites.CSS(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500, %s", err)
			log.Println(err)
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
