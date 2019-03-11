package controllers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"time"
)

// SignIn godoc
// @title Sign in
// @summary Grants client access
// @description This method logs user in and sets cookie
// @tags auth
// @accept json
// @produce json
// @param credentials body models.SignInData true "Credentials"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /session [post]
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

// SignOut godoc
// @title Sign out
// @summary Logs user out
// @description This method logs user out and deletes cookie
// @tags auth
// @produce json
// @param credentials body models.SignInData true "Credentials"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /session [delete]
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

// IsAuth godoc
// @title Check session
// @summary Checks user session
// @description This method checks whether user is signed in or signed out
// @tags auth
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /session [get]
func IsAuth(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusOK, "signed in")
	} else {
		models.SendMessage(w, http.StatusUnauthorized, "signed out")
	}
}
