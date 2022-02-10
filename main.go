package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Create two separate HTTP servers to demonstrate proxy as backends to receive
	// requests from rproxy.
	go runInternalDemoServer(":8888")
	go runInternalDemoServer(":9999")

	// This is just to make using go test in main slightly less janky.
	doProxy()
}

func runInternalDemoServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello from %s", addr)
	})
	fmt.Printf("Listening on http://%s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func doProxy() {
	configServer()
	runServer()
}

func configServer() {
	targetA, err := url.Parse("http://localhost:8888")
	if err != nil {
		log.Fatalf("%v", err)
	}

	targetB, err := url.Parse("http://localhost:9999")
	if err != nil {
		log.Fatalf("%v", err)
	}

	rpA := httputil.NewSingleHostReverseProxy(targetA)
	rpB := httputil.NewSingleHostReverseProxy(targetB)

	svcA := svcA{rp: rpA}
	svcB := svcB{rp: rpB}

	// This configures http.DefaultServeMux to handle all requrests using httputil.ReverseProxy server.
	// This function could alternatively return an http.ServeMux instead of operating on the default.
	//
	// the call to forceLogReq creates a 'middleware' function that chains http.Handler's in succession.
	http.Handle("/svca", forceLogReq(svcA))
	http.Handle("/svcb", forceLogReq(svcB))
}

func runServer() {
	// Listens on 0.0.0.0:8080 and ::8080
	addr := ":8080"
	fmt.Printf("Listening on http://%s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func forceLogReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Docs say this is expensive is likely not desired for production.
		bytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("%v", err)
		} else {
			log.Printf("%s\n", bytes)
		}
		next.ServeHTTP(w, r)
	})
}

// Normally this would likely go in a sub package and be exported.
type svcA struct {
	rp *httputil.ReverseProxy
}

func (s svcA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Just pass the request on to the configured http.ReverseProxy
	// Additional logic or request modification could be handled here
	s.rp.ServeHTTP(w, r)
}

// Normally this would likely go in a sub package and be exported.
type svcB struct {
	rp *httputil.ReverseProxy
}

func (s svcB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Just pass the request on to the configured http.ReverseProxy
	// Additional logic or request modification could be handled here
	s.rp.ServeHTTP(w, r)
}
