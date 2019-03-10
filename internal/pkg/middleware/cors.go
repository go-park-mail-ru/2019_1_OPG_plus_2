package middleware

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

type CorsData struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	MaxAge           int
	AllowCredentials bool
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		val, ok := req.Header["Origin"]
		if ok {
			res.Header().Set("Access-Control-Allow-Origin", val[0])
		}
		next.ServeHTTP(res, req)
	})
}

func ApplyJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(res, req)
	})
}
