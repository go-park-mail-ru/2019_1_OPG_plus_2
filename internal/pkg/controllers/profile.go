package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var MB = 1 << 20

var userStorage = models.NewUserProfileStorage()

func init() {
	_ = userStorage.Set(1, &models.UserProfile{
		ID:        1,
		Score:     228,
		AvatarUrl: "<user1_avatar_url>",
	})

	_ = userStorage.Set(2, &models.UserProfile{
		ID:        2,
		Score:     1337,
		AvatarUrl: "<user2_avatar_url>",
	})
}

//TODO: в ScoreBoard пагинация по limit offset

// CreateProfile godoc
// @title Create profile
// @summary Registers user
// @description This method creates records about new user in auth-bd and profile-db and then sends cookie to user in order to identify
// @accept json
// @produce json
// @param profile_data body models.UserData true "User profile data"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /profile [post]
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var userData models.UserData
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		models.SendMessage(w, http.StatusBadRequest, "JSON parsing error")
		return
	}

	jwtData, err := auth.CreateUser(userData)
	if err != nil {
		models.SendMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))

	var profile models.UserProfile
	profile.ID = jwtData.Id
	profile.Username = jwtData.Username
	profile.Email = jwtData.Email
	_ = userStorage.Set(jwtData.Id, &profile)
	models.SendMessage(w, http.StatusOK, "Profile successfully created")
}

// GetProfile godoc
// @title Get profile
// @summary Produces user profile info
// @description This method provides client with profile data, matching required ID
// @accept json
// @produce json
// @param id path int true "Profile ID"
// @success 200 {object} models.UserProfile
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 404 {object} models.SuccessOrErrorMessage
// @router /profile/{id} [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)

	if pathVariables == nil {
		models.SendMessage(w, http.StatusBadRequest, "Bad query")
		return
	}
	id, ok := pathVariables["id"]
	if !ok {
		models.SendMessage(w, http.StatusBadRequest, "Bad query")
		return
	}

	intId, _ := strconv.ParseInt(id, 10, 64)
	profile, err := userStorage.Get(int(intId))
	if err != nil {
		models.SendMessage(w, http.StatusNotFound, "User not found")
		return
	}
	msg, _ := json.Marshal(profile)

	_, _ = fmt.Fprintln(w, string(msg))
}

// UpdateProfile godoc
// @title Update profile
// @summary Updates client's profile
// @description This method updates info in profile and auth-db record of user, who is making a query
// @accept json
// @produce json
// @param profile_data body models.UserData true "User new profile data"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /profile [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// UpdateEmail
	// UpdateUsername
	// whatever else...

	var newProfile models.UserData
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		models.SendMessage(w, http.StatusBadRequest, "JSON parsing error")
		return
	}

	id := jwtData(r).Id
	user, _ := userStorage.Get(id)
	user.Username = newProfile.Username
	user.Email = newProfile.Email

	//auth.UpdateUser(newProfile)

	models.SendMessage(w, http.StatusOK, "Profile successfully updated")
}

// DeleteProfile godoc
// @title Delete profile
// @summary Deletes profile and user of client
// @description This method deletes all information about user, making a query, including profile, game stats and authorization info
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /profile [delete]
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := jwtData(r).Id

	err := userStorage.Delete(id)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "Error while deleting profile")
	}

	models.SendMessage(w, http.StatusOK, "Profile "+string(id)+" deleted successfully")
}

// TODO: implement FileStorage interface to inline file operations

// UploadAvatar godoc
// @title Upload new avatar
// @summary Saves new avatar image of client's user
// @description This method saves avatar image in server storage and sets it as clients user avatar
// @accept png
// @accept jpeg
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @router /upload_avatar [post]
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(int64(5 * MB))
	file, _, err := r.FormFile("avatar")
	if err != nil {
		models.SendMessage(w, http.StatusBadRequest, "Error reading file")
		return
	}
	defer file.Close()

	data := bytes.NewBuffer(nil)
	_, err = io.Copy(data, file)
	if err != nil {
		panic(err)
	}

	id := jwtData(r).Id
	fmt.Println(id)
	user, _ := userStorage.Get(id)

	err = ioutil.WriteFile(`~/colors_static/`+strconv.FormatInt(int64(id), 10)+`.png`, data.Bytes(), 0666)
	if err != nil {
		panic(err)
	}

	user.AvatarUrl = "/static/" + strconv.FormatInt(int64(id), 10) + `.png`

	models.SendMessage(w, http.StatusOK, user.AvatarUrl)
}

//сервисный метод чтобы понимать что творится в хранилище юзеров
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	message, _ := json.Marshal(userStorage.Data)
	_, _ = fmt.Fprintln(w, string(message))
}
