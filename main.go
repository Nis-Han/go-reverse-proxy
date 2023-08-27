package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

type Route struct {
	Path   string `json:"path"`
	Target string `json:"target"`
}

var routes []Route

func main() {
	loadRoutesFromJSON("routes.json")

	r := mux.NewRouter()

	for _, route := range routes {
		fmt.Println(route.Path, route.Target)
		// addReverseProxyRoute(r, route.Path, route.Target)
	}

	port := 8080
	fmt.Printf("Reverse proxy server is running on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadRoutesFromJSON(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&routes)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
}

func addReverseProxyRoute(r *mux.Router, path string, targetURL string) {
	target, _ := url.Parse(targetURL)

	r.PathPrefix(path).Handler(http.StripPrefix(path, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		proxy := httputil.NewSingleHostReverseProxy(target)
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.Header.Set("X-Real-IP", req.RemoteAddr)
		proxy.ServeHTTP(w, req)
	})))
}
