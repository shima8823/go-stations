package middleware

import (
	"context"
	"github.com/mileusna/useragent"
	"net/http"
)

type osNameKeyType string

const OSNameKey osNameKeyType = "osname"

func OSParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), OSNameKey, useragent.Parse(r.UserAgent()).OS)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
