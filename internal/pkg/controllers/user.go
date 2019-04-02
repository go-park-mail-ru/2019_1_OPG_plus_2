package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	a "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/adapters"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func NewUserHandlers() *UserHandlers {
	return &UserHandlers{}
}

type UserHandlers struct{}

// CreateUser godoc
// @title Create User
// @summary Registers User
// @description This method creates records about new User in Auth-bd and User-db and then sends cookie to User in order to identify
// @tags User
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User data"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.IncorrectFieldsAnswer
// @failure 405 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /User [post]
func (*UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	jwtData, err, fields := a.GetStorages().User.CreateUser(signUpData)
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
// @title Get User
// @summary Produces User profile info
// @description This method provides client with User data, matching required ID
// @tags User
// @accept json
// @produce json
// @param id path int false "ProfileData ID, if none, returned logged in User"
// @success 200 {object} models.UserDataAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 404 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /User/{id} [get]
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
		return
	}

	models.Send(w, http.StatusOK, models.GetUserDataAnswer(userData))
}

// UpdateUser godoc
// @title Update User
// @summary Updates client's User
// @description This method updates info in profile and Auth-db record of User, who is making a query
// @tags User
// @accept json
// @produce json
// @param profile_data body models.UpdateUserData true "User new profile data"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /User [put]
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
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.Send(w, http.StatusOK, models.UserUpdatedAnswer)
}

// RemoveUser godoc
// @title Delete User
// @summary Deletes User and User of client
// @description This method deletes all information about User, making a query, including profile, game stats and authorization info
// @tags User
// @produce json
// @param remove_data body models.RemoveUserData true "Info required to remove current User"
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 401 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /User [delete]
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
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData(r), 0))
	models.Send(w, http.StatusOK, models.UserRemovedAnswer)
}
