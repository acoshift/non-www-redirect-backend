package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

var port = flag.Int("port", 8080, "Port to server non www redirect backend")

func main() {
	flag.Parse()
	http.Handle("/", nonWWWRedirect()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "non www redirect backend - 404")
	})))
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start http server: %s\n", err)
		os.Exit(1)
	}
}

func nonWWWRedirect() middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host := r.Header.Get(header.XForwardedHost)
			if len(host) == 0 {
				host = r.Host
			}
			host = strings.TrimPrefix(host, "www.")
			if len(host) < len(r.Host) {
				http.Redirect(w, r, scheme(r)+"://"+host+r.RequestURI, http.StatusMovedPermanently)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func isTLS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if r.Header.Get(header.XForwardedProto) == "https" {
		return true
	}
	return false
}

func scheme(r *http.Request) string {
	if isTLS(r) {
		return "https"
	}
	return "http"
}
