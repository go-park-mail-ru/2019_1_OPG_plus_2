package middleware

import (
	"fmt"
	"net/http"
	"reflect"
)

func ValueOfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(reflect.ValueOf(next))
		next.ServeHTTP(res, req)
	})
}
