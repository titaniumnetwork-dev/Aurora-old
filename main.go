package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", server)
	http.ListenAndServe(":8080", nil)
}

func server(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}