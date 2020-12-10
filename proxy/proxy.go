package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"fmt"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
        tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 10 * time.Second,
	}

	client := &http.Client{Transport: tr}

	url := strings.Join(strings.Split(r.URL.String(), "?")[1:], "")
	if url == "" {
		fmt.Fprintf(w, "Welcome to %s proxy!", "Aurora")
                return
	}

	if strings.HasPrefix(url, "//") {
                url = "http:" + url
        } else if !strings.HasPrefix(url, "http") {
                url = "https://" + url
        }

	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	for key, val := range resp.Header {
		w.Header().Set(key, strings.Join(val, ", "))
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
