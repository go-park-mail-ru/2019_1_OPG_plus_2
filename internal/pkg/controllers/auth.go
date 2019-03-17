package controllers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"time"
)

// IsAuth godoc
// @title Check session
// @summary Checks user session
// @description This method checks whether user is signed in or signed out
// @tags auth
// @produce json
// @success 200 {object} models.AnswerMessage
// @failure 401 {object} models.AnswerMessage
// @router /session [get]
func IsAuth(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusOK, "signed in")
	} else {
		models.SendMessage(w, http.StatusUnauthorized, "signed out")
	}
}

// SignIn godoc
// @title Sign in
// @summary Grants client access
// @description This method logs user in and sets cookie
// @tags auth
// @accept json
// @produce json
// @param credentials body models.SignInData true "Credentials"
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessageWithFields
// @failure 405 {object} models.AnswerMessage
// @failure 500 {object} models.AnswerMessage
// @router /session [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusMethodNotAllowed, "already signed in")
		return
	}

	signInData := models.SignInData{}
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "incorrect JSON")
		return
	}
	defer r.Body.Close()

	jwtData, err, fields := auth.SignIn(signInData)
	if err != nil {
		if fields != nil {
			models.SendMessageWithFields(w, http.StatusBadRequest, err.Error(), fields)
			return
		}
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
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
// @success 200 {object} models.AnswerMessage
// @failure 405 {object} models.AnswerMessage
// @router /session [delete]
func SignOut(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.SendMessage(w, http.StatusMethodNotAllowed, "already signed out")
		return
	}

	jwtCookie, _ := r.Cookie(auth.CookieName)
	jwtCookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, jwtCookie)
	models.SendMessage(w, http.StatusOK, "signed out")
}

// UpdatePassword godoc
// @title Update password
// @summary Updates user password
// @description This method updates users password, requiring password and confirmation. User data is pulled from jwt-token
// @tags auth
// @accepts json
// @produce json
// @param update_data body models.UpdatePasswordData true "New password info"
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessageWithFields
// @failure 401 {object} models.AnswerMessage
// @failure 500 {object} models.AnswerMessage
// @router /password [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.SendMessage(w, http.StatusUnauthorized, "not signed in")
		return
	}

	updateData := models.UpdatePasswordData{}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "incorrect JSON")
		return
	}
	defer r.Body.Close()

	err, fields := auth.UpdatePassword(jwtData(r).Id, updateData)
	if err != nil {
		if fields != nil {
			models.SendMessageWithFields(w, http.StatusBadRequest, err.Error(), fields)
			return
		}
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.SendMessage(w, http.StatusOK, "password updated")
}
