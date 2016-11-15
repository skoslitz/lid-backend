package main

import (
	// "log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messenger := color.New(color.FgYellow).PrintfFunc()
		start := time.Now()

		inner.ServeHTTP(w, r)

		defer messenger("\n+--------------------------------------------------------------+\n")

		// log.Printf
		messenger(
			"%s  %s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
			// name,

		)
	})
}
