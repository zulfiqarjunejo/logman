package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zulfiqarjunejo/logs-management-system/clients"
	"github.com/zulfiqarjunejo/logs-management-system/types"
	"golang.org/x/crypto/bcrypt"
)

func CreateCheckApiKey(clientModel clients.ClientModel) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			clientId := r.Header.Get("x-client-id")
			apiKey := r.Header.Get("x-api-key")

			if clientId == "" || apiKey == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			client, err := clientModel.FindClientById(clientId)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(client.ApiKey), []byte(apiKey))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			const clientContextKey types.ContextKey = 0
			contextWithClient := context.WithValue(r.Context(), clientContextKey, client)

			next.ServeHTTP(w, r.WithContext(contextWithClient))
		}

		return http.HandlerFunc(fn)
	}
}

func PrintRouteInfo(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL
		method := r.Method

		fmt.Printf("%s %s\n", method, path)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
