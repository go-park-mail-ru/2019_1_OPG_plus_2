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
		"http://localhost:8002",
		"http://localhost:8003",
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
		origin, hasOrigin := req.Header["Origin"]
		if hasOrigin {
			found := false
			for _, allowed := range corsData.AllowOrigins {
				if origin[0] == allowed {
					found = true
					break
				}
			}
			if found {
				res.Header().Set("Access-Control-Allow-Origin", origin[0])
				res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsData.AllowCredentials))
			} else {
				fmt.Println("Origin " + origin[0] + " wasn't found!")
			}
		}

		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Methods", strings.Join(corsData.AllowMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(corsData.AllowHeaders, ", "))
			res.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsData.MaxAge))
			return
		}

		next.ServeHTTP(res, req)
	})
}
