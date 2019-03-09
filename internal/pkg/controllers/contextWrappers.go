package controllers

import "net/http"

func isAuth(r *http.Request) bool {
	return true
}

func userGuid(r *http.Request) string {
	return "1"
}
