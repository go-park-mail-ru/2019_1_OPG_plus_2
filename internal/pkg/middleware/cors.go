package middleware

import (
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
        "https://api.colors.hackallcode.ru",
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
        val, ok := req.Header["Origin"]
        if ok {
            res.Header().Set("Access-Control-Allow-Origin", val[0])
            res.Header().Set("Access-Control-Allow-Credentials", "true")
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

func ApplyJsonContentType(next http.Handler) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.Header().Set("Content-Type", "application/json; charset=utf-8")
        next.ServeHTTP(res, req)
    })
}
