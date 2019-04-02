package middleware

import (
	"net/http"
)

func ApplyJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(res, req)
	})
}
