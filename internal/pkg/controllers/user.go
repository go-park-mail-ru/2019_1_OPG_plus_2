package controllers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/user"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

var storageAdapter = user.NewStorageAdapter()

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
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessage
// @failure 401 {object} models.AnswerMessage
// @router /user [post]
func (*UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if isAuth(r) {
		models.SendMessage(w, http.StatusBadRequest, "already signed in")
		return
	}

	signUpData := models.SingUpData{}
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "incorrect JSON")
		return
	}
	defer r.Body.Close()

	jwtData, err := storageAdapter.CreateUser(signUpData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.SendMessage(w, http.StatusOK, "signed up")
}

// GetUser godoc
// @title Get user
// @summary Produces user profile info
// @description This method provides client with user data, matching required ID
// @tags user
// @accept json
// @produce json
// @param id path int true "Profile ID"
// @success 200 {object} models.UserData
// @failure 400 {object} models.AnswerMessage
// @failure 404 {object} models.AnswerMessage
// @router /user/{id} [get]
func (*UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	var id int64
	var err error

	vars := mux.Vars(r)
	pathId, ok := vars["id"]
	if ok {
		id, err = strconv.ParseInt(pathId, 10, 64)
		if err != nil {
			models.SendMessage(w, http.StatusBadRequest, "incorrect id in query")
			return
		}
	} else {
		if !isAuth(r) {
			models.SendMessage(w, http.StatusBadRequest, "no id in query")
			return
		}
		id = jwtData(r).Id
	}

	userData, err := storageAdapter.GetUser(id)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.SendMessageWithData(w, http.StatusOK, "user found", userData)
}

// UpdateUser godoc
// @title Update user
// @summary Updates client's user
// @description This method updates info in profile and auth-db record of user, who is making a query
// @tags user
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User new profile data"
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessage
// @failure 401 {object} models.AnswerMessage
// @router /user [put]
func (*UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.SendMessage(w, http.StatusUnauthorized, "not signed in")
		return
	}

	var updateData models.UpdateUserData
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "incorrect JSON")
		return
	}
	defer r.Body.Close()

	jwtData, err := storageAdapter.UpdateUser(jwtData(r).Id, updateData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
	models.SendMessage(w, http.StatusOK, "user updated")
}

// RemoveUser godoc
// @title Delete user
// @summary Deletes user and user of client
// @description This method deletes all information about user, making a query, including profile, game stats and authorization info
// @tags user
// @produce json
// @success 200 {object} models.AnswerMessage
// @failure 500 {object} models.AnswerMessage
// @router /user [delete]
func (*UserHandlers) RemoveUser(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.SendMessage(w, http.StatusUnauthorized, "not signed in")
		return
	}

	var removeData models.RemoveUserData
	err := json.NewDecoder(r.Body).Decode(&removeData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "incorrect JSON")
		return
	}
	defer r.Body.Close()

	err = storageAdapter.RemoveUser(jwtData(r).Id, removeData)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData(r), 0))
	models.SendMessage(w, http.StatusOK, "user removed")
}