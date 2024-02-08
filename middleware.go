package main

import (
	"fmt"
	"net/http"
)

func RouteInfoMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL
		method := r.Method

		fmt.Printf("%s %s\n", method, path)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
