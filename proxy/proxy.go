package proxy

import (
	"io"
	"log"
	"net/http"
)

// Server used for proxy
func Server(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
