package main

import (
    "fmt"
    "hash/fnv"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
    "net/url"
)

var backends = []string{
    "http://192.0.0.80:8009",
    "http://192.0.1.80:8009",
}

// Hash function to mimic IP Hash logic
func hashIP(ip string) int {
    h := fnv.New32a()
    h.Write([]byte(ip))
    return int(h.Sum32()) % len(backends)
}

// Reverse proxy handler
func proxyHandler(w http.ResponseWriter, r *http.Request) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        http.Error(w, "Invalid IP address", http.StatusInternalServerError)
        return
    }

    backendIndex := hashIP(ip)
    target, err := url.Parse(backends[backendIndex])
    if err != nil {
        http.Error(w, "Bad backend URL", http.StatusInternalServerError)
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(target)

    // Rewrite request's Host header
    r.Host = target.Host
    proxy.ServeHTTP(w, r)
}

func main() {
    fmt.Println("Go API Load Balancer running on :3309")
    http.HandleFunc("/", proxyHandler)

    err := http.ListenAndServe(":3309", nil)
    if err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
