package controllers

import (
	"net/http"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

func IsAuth(r *http.Request) bool {
	return r.Context().Value("isAuth").(bool)
}

func JwtData(r *http.Request) models.JwtData {
	return r.Context().Value("jwtData").(models.JwtData)
}
