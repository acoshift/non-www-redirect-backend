package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var port = flag.Int("port", 8080, "Port to serve non www redirect backend")

func main() {
	flag.Parse()

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host

		redirectHost := strings.TrimPrefix(host, "www.")
		if len(redirectHost) < len(host) {
			http.Redirect(w, r, scheme(r)+"://"+redirectHost+r.RequestURI, http.StatusMovedPermanently)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "non www redirect backend - 404")
	}))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "could not start http server: %s\n", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "could not shutdown http server: %s\n", err)
	}
}

func isTLS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if r.Header.Get("X-Forwarded-Proto") == "https" {
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
