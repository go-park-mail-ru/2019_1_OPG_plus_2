package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/user"
)

type UserHandlers struct {
	adapter user.IUser
}

func NewUserHandlers(adapter user.IUser) *UserHandlers {
	return &UserHandlers{adapter: adapter}
}

// CreateUser godoc
// @title Create user
// @summary Registers user
// @description This method creates records about new user in auth-bd and user-db and then sends cookie to user in order to identify
// @tags user
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User data"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.IncorrectFieldsAnswer
// @failure 405 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user [post]
func (handlers *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.Send(w, http.StatusMethodNotAllowed, models.AlreadySignedInAnswer)
		return
	}

	signUpData := models.SingUpData{}
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.IncorrectJsonAnswer)
		return
	}
	defer r.Body.Close()

	jwtData, err, fields := handlers.adapter.CreateUser(signUpData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
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
// @param id path int false "ProfileData ID, if none, returned logged in user"
// @success 200 {object} models.UserDataAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 404 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user/{id} [get]
func (handlers *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
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

	userData, err := handlers.adapter.GetUser(id)
	if err != nil {
		if err == models.NotFound {
			models.Send(w, http.StatusNotFound, models.UserNotFoundAnswer)
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		return
	}

	models.Send(w, http.StatusOK, models.GetUserDataAnswer(userData))
}

// UpdateUser godoc
// @title Update user
// @summary Updates client's user
// @description This method updates info in profile and auth-db record of user, who is making a query
// @tags user
// @accept json
// @produce json
// @param profile_data body models.UpdateUserData true "User new profile data"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user [put]
func (handlers *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	jwtData, err, fields := handlers.adapter.UpdateUser(jwtData(r).Id, updateData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.Send(w, http.StatusOK, models.UserUpdatedAnswer)
}

// RemoveUser godoc
// @title Delete user
// @summary Deletes user and user of client
// @description This method deletes all information about user, making a query, including profile, game stats and authorization info
// @tags user
// @produce json
// @param remove_data body models.RemoveUserData true "Info required to remove current user"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /user [delete]
func (handlers *UserHandlers) RemoveUser(w http.ResponseWriter, r *http.Request) {
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

	err, fields := handlers.adapter.RemoveUser(jwtData(r).Id, removeData)
	if err != nil {
		if fields != nil {
			models.Send(w, http.StatusBadRequest, models.GetIncorrectFieldsAnswer(fields))
			return
		}
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData(r), 0))
	models.Send(w, http.StatusOK, models.UserRemovedAnswer)
}
