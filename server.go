package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/acoshift/middleware"
)

var port = flag.Int("port", 8080, "Port to server non www redirect backend")

func main() {
	flag.Parse()
	http.Handle("/", middleware.NonWWWRedirect()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
