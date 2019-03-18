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
// @success 200 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @router /session [get]
func IsAuth(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.Send(w, http.StatusOK, models.SignedInAnswer)
	} else {
		models.Send(w, http.StatusUnauthorized, models.SignedOutAnswer)
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
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.IncorrectFieldsAnswer
// @failure 405 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /session [post]
func SignIn(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.Send(w, http.StatusMethodNotAllowed, models.AlreadySignedInAnswer)
		return
	}

	signInData := models.SignInData{}
	err := json.NewDecoder(r.Body).Decode(&signInData)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.IncorrectJsonAnswer)
		return
	}
	defer r.Body.Close()

	jwtData, err, fields := auth.SignIn(signInData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.NewIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.MessageAnswer{Status: 500, Message: err.Error()})
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.Send(w, http.StatusOK, models.SignedInAnswer)
}

// SignOut godoc
// @title Sign out
// @summary Logs user out
// @description This method logs user out and deletes cookie
// @tags auth
// @produce json
// @success 200 {object} models.MessageAnswer
// @failure 405 {object} models.MessageAnswer
// @router /session [delete]
func SignOut(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.Send(w, http.StatusMethodNotAllowed, models.AlreadySignedOutAnswer)
		return
	}

	jwtCookie, _ := r.Cookie(auth.CookieName)
	jwtCookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, jwtCookie)
	models.Send(w, http.StatusOK, models.SignedOutAnswer)
}

// UpdatePassword godoc
// @title Update password
// @summary Updates user password
// @description This method updates users password, requiring password and confirmation. User data is pulled from jwt-token
// @tags auth
// @accepts json
// @produce json
// @param update_data body models.UpdatePasswordData true "New password info"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.IncorrectFieldsAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /password [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.Send(w, http.StatusUnauthorized, models.NotSignedInAnswer)
		return
	}

	updateData := models.UpdatePasswordData{}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.IncorrectJsonAnswer)
		return
	}
	defer r.Body.Close()

	err, fields := auth.UpdatePassword(jwtData(r).Id, updateData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.NewIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.MessageAnswer{Status: 500, Message: err.Error()})
		return
	}

	models.Send(w, http.StatusOK, models.PasswordUpdatedAnswer)
}
