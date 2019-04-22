package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"2019_1_OPG_plus_2/internal/pkg/tsLogger"

	a "2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{}
}

type AuthHandlers struct{}

// IsAuth godoc
// @title Check session
// @summary Checks User session
// @description This method checks whether User is signed in or signed out
// @tags Auth
// @produce json
// @success 200 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @router /session [get]
func (*AuthHandlers) IsAuth(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.Send(w, http.StatusOK, models.SignedInAnswer)
	} else {
		models.Send(w, http.StatusUnauthorized, models.SignedOutAnswer)
	}
}

// SignIn godoc
// @title Sign in
// @summary Grants client access
// @description This method logs User in and sets cookie
// @tags Auth
// @accept json
// @produce json
// @param credentials body models.SignInData true "Credentials"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.IncorrectFieldsAnswer
// @failure 405 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /session [post]
func (*AuthHandlers) SignIn(w http.ResponseWriter, r *http.Request) {
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

	jwtData, err, fields := a.GetStorages().Auth.SignIn(signInData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		tsLogger.LogErr(fmt.Sprintf("DEV ERR: %q ==> %v", r.RequestURI, err))
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.Send(w, http.StatusOK, models.SignedInAnswer)
}

// SignOut godoc
// @title Sign out
// @summary Logs User out
// @description This method logs User out and deletes cookie
// @tags Auth
// @produce json
// @success 200 {object} models.MessageAnswer
// @failure 405 {object} models.MessageAnswer
// @router /session [delete]
func (*AuthHandlers) SignOut(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.Send(w, http.StatusMethodNotAllowed, models.AlreadySignedOutAnswer)
		return
	}

	jwtCookie := auth.CreateAuthCookie(models.JwtData{}, -1)
	http.SetCookie(w, jwtCookie)
	models.Send(w, http.StatusOK, models.SignedOutAnswer)
}

// UpdatePassword godoc
// @title Update password
// @summary Updates User password
// @description This method updates users password, requiring password and confirmation. User data is pulled from jwt-token
// @tags Auth
// @accepts json
// @produce json
// @param update_data body models.UpdatePasswordData true "New password info"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.IncorrectFieldsAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /password [put]
func (*AuthHandlers) UpdatePassword(w http.ResponseWriter, r *http.Request) {
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

	err, fields := a.GetStorages().Auth.UpdatePassword(jwtData(r).Id, updateData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		tsLogger.LogErr(fmt.Sprintf("DEV ERR: %q ==> %v", r.RequestURI, err))
		return
	}

	models.Send(w, http.StatusOK, models.PasswordUpdatedAnswer)
}
