package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedUserID := os.Getenv("BASIC_AUTH_USER_ID")
		expectedPassword := os.Getenv("BASIC_AUTH_PASSWORD")
		if expectedUserID == "" || expectedPassword == "" {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok || username != expectedUserID || password != expectedPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
