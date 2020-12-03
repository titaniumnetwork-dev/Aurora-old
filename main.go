package main

import (
	"AuroraProxy/proxy"
	"net/http"
)

func main() {
	http.HandleFunc("/", proxy.Server)
	http.ListenAndServe(":8080", nil)
}
