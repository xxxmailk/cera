package main

import (
	proxy "cera/http"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/", &proxy.Proxy{IP: "192.168.56.79", Port: "3000"})
	panic(http.ListenAndServe("0.0.0.0:8080", nil))
}
