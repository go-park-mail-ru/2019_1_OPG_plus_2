package controllers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"time"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusBadRequest, "already signed in")
		return
	}

	signInData := models.SignInData{}
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "can't parse request body")
		return
	}
	defer r.Body.Close()

	jwtData, err := auth.CheckLoginPass(signInData)
	if err != nil {
		models.SendMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.SendMessage(w, http.StatusOK, "signed in")
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.SendMessage(w, http.StatusUnauthorized, "already signed out")
		return
	}

	jwtCookie, errNoCookie := r.Cookie(auth.CookieName)
	if errNoCookie != nil {
		models.SendMessage(w, http.StatusUnauthorized, "already signed out")
		return
	}

	jwtCookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, jwtCookie)
	models.SendMessage(w, http.StatusOK, "signed out")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusBadRequest, "already signed in")
		return
	}

	userData := models.UserData{}
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "can't parse request body")
		return
	}
	defer r.Body.Close()

	jwtData, err := auth.CreateUser(userData)
	if err != nil {
		models.SendMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.SendMessage(w, http.StatusOK, "signed up")
}

func IsAuth(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusOK, "signed in")
	} else {
		models.SendMessage(w, http.StatusUnauthorized, "signed out")
	}
}
