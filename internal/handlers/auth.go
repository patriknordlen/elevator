package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

type key int
const userKey key = iota

func MockUser(next http.HandlerFunc, user string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r.WithContext(context.WithValue(r.Context(), userKey, user)))
	}
}

func RequireToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) == 2 && authHeader[0] == "Bearer" {
			idToken := authHeader[1]
			parsedToken, err := idtoken.Validate(r.Context(), idToken, "")

			if err != nil {
				log.Println("Error: ", err)
				ReturnUnauthorized(w)
			} else {
				next(w, r.WithContext(context.WithValue(r.Context(), userKey, parsedToken.Claims["email"])))
			}
		} else {
			ReturnUnauthorized(w)
		}
	}
}
