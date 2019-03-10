package controllers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err := fmt.Fprintf(w, "Site of OPG+2!")
	if err != nil {
		fmt.Println(err)
	}
}

func IndexApiHandler(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusOK, "Hello, "+jwtData(r).Username+"!")
	} else {
		models.SendMessage(w, http.StatusOK, "I don't know about you, but hello!")
	}
}
