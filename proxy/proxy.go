package proxy

import (
	"fmt"
	"net/http"
)

// Server to be used for the proxy.
func Server(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
