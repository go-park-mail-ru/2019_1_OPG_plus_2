package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Middleware func(next http.Handler) http.Handler

type CorsData struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	MaxAge           int
	AllowCredentials bool
}

var corsData = CorsData{
	AllowOrigins: []string{
		"https://colors.hackallcode.ru",
		"http://localhost:8001",
		"http://127.0.0.1:8001",
	},
	AllowMethods: []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
	},
	AllowHeaders: []string{
		"Content-Type",
	},
	MaxAge:           88500,
	AllowCredentials: true,
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			res.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsData.MaxAge))
			return
		}

		origin, ok := req.Header["Origin"]
		if !ok {
			return
		}
		found := false
		for _, v := range corsData.AllowOrigins {
			if origin[0] == v {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Origin " + origin[0] + " wasn't found!")
			return
		}

		res.Header().Set("Access-Control-Allow-Origin", origin[0])
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(res, req)
	})
}
