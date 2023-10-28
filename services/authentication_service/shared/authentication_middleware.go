package shared

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"net/http"
	"strings"
)

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	bearerAndToken := strings.Split(authHeader, " ")

	if len(bearerAndToken) != 2 {
		return ""
	}

	return bearerAndToken[1]
}

func AuthenticationMiddleware(authClient *auth.Client) func(w http.ResponseWriter, r *http.Request, next func(), abort func(code int)) {
	return func(w http.ResponseWriter, r *http.Request, next func(), abort func(code int)) {

		token := extractToken(r)

		idToken, err := authClient.VerifyIDToken(r.Context(), token)
		if err != nil {
			abort(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", idToken.UID)
		newRequest := r.WithContext(ctx)
		*r = *newRequest
	}
}
