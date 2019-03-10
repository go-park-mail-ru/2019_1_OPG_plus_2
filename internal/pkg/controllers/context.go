package controllers

import (
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
    "net/http"
)

func isAuth(r *http.Request) bool {
    return r.Context().Value("isAuth").(bool)
}

func jwtData(r *http.Request) models.JwtData {
    return r.Context().Value("jwtData").(models.JwtData)
}
