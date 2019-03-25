package controllers

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	models.Send(w, http.StatusOK, models.GetSuccessAnswer("Backend of OPG+2 project 'Colors'!"))
}

func IndexApiHandler(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.Send(w, http.StatusOK, models.GetSuccessAnswer("Hello, "+jwtData(r).Username+"!"))
	} else {
		models.Send(w, http.StatusOK, models.GetSuccessAnswer("I don't know about you, but hello!"))
	}
}
