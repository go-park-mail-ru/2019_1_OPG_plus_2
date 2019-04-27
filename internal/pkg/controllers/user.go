package controllers

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	a "2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/auth"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

func NewUserHandlers() *UserHandlers {
	return &UserHandlers{}
}

type UserHandlers struct{}

// CreateUser godoc
// @title Create user
// @summary Registers user
// @description This method creates records about new user in Auth-bd and user-db and then sends cookie to user in order to identify
// @tags user
// @accept json
// @produce json
// @param profile_data body models.SignUpData true "user data"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.IncorrectFieldsAnswer
// @failure 405 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user [post]
func (*UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.Send(w, http.StatusMethodNotAllowed, models.AlreadySignedInAnswer)
		return
	}

	signUpData := models.SignUpData{}
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.IncorrectJsonAnswer)
		return
	}
	defer r.Body.Close()

	jwtData, err, fields := a.GetStorages().User.CreateUser(signUpData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		tsLogger.LogErr("DEV ERR: %q ==> %v", r.RequestURI, err)
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.Send(w, http.StatusOK, models.SignedUpAnswer)
}

// GetUser godoc
// @title Get user
// @summary Produces user profile info
// @description This method provides client with user data, matching required ID
// @tags user
// @accept json
// @produce json
// @param id path int false "users ID, if none, returned logged in user"
// @success 200 {object} models.UserDataAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 404 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user/{id} [get]
func (*UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	var id int64
	var err error

	vars := mux.Vars(r)
	pathId, ok := vars["id"]
	if ok {
		id, err = strconv.ParseInt(pathId, 10, 64)
		if err != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer([]string{"id"}))
			return
		}
	} else {
		if !isAuth(r) {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer([]string{"id"}))
			return
		}
		id = jwtData(r).Id
	}

	userData, err := a.GetStorages().User.GetUser(id)
	if err != nil {
		if err == models.NotFound {
			models.Send(w, http.StatusNotFound, models.UserNotFoundAnswer)
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		tsLogger.LogErr("DEV ERR: %q ==> %v", r.RequestURI, err)
		return
	}

	models.Send(w, http.StatusOK, models.GetUserDataAnswer(userData))
}

// UpdateUser godoc
// @title Update user
// @summary Updates client's user
// @description This method updates info in profile and Auth-db record of user, who is making a query
// @tags user
// @accept json
// @produce json
// @param profile_data body models.UpdateUserData true "user new profile data"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user [put]
func (*UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.Send(w, http.StatusUnauthorized, models.NotSignedInAnswer)
		return
	}

	var updateData models.UpdateUserData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.IncorrectJsonAnswer)
		return
	}
	defer r.Body.Close()

	jwtData, err, fields := a.GetStorages().User.UpdateUser(jwtData(r).Id, updateData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		tsLogger.LogErr("DEV ERR: %q ==> %v", r.RequestURI, err)
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.Send(w, http.StatusOK, models.UserUpdatedAnswer)
}

// RemoveUser godoc
// @title Delete user
// @summary Deletes user and profile of client
// @description This method deletes all information about user, making a query, including profile, game stats and authorization info
// @tags user
// @produce json
// @param remove_data body models.RemoveUserData true "Info required to remove current user"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user [delete]
func (*UserHandlers) RemoveUser(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.Send(w, http.StatusUnauthorized, models.NotSignedInAnswer)
		return
	}

	var removeData models.RemoveUserData
	err := json.NewDecoder(r.Body).Decode(&removeData)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.IncorrectJsonAnswer)
		return
	}
	defer r.Body.Close()

	err, fields := a.GetStorages().User.RemoveUser(jwtData(r).Id, removeData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		tsLogger.LogErr("DEV ERR: %q ==> %v", r.RequestURI, err)
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData(r), 0))
	models.Send(w, http.StatusOK, models.UserRemovedAnswer)
}
