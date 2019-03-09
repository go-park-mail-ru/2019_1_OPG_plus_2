package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var MB = 1 << 20

var userStorage = models.NewUserProfileStorage()

func init() {
	_ = userStorage.Set("1", &models.UserProfile{
		Username:  "user1",
		Email:     "user1@example.com",
		Score:     228,
		AvatarUrl: "<user1_avatar_url>",
	})

	_ = userStorage.Set("2", &models.UserProfile{
		Username:  "user2",
		Email:     "user2@example.com",
		Score:     1337,
		AvatarUrl: "<user2_avatar_url>",
	})
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.UserProfile
	var message models.SuccessOrErrorMessage
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		message.Send(w, http.StatusBadRequest, "JSON parsing error")
		return
	}
	_ = userStorage.Set(strconv.FormatInt(int64(len(userStorage.Data)+1), 10), &profile)
	message.Send(w, http.StatusOK, "Profile successfully created")
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	var message models.SuccessOrErrorMessage
	if pathVariables == nil {
		message.Send(w, http.StatusBadRequest, "Bad query")
		return
	}
	id, ok := pathVariables["id"]
	if !ok {
		message.Send(w, http.StatusBadRequest, "Bad query")
		return
	}

	profile, err := userStorage.Get(id)
	if err != nil {
		message.Send(w, http.StatusNotFound, "User not found")
		return
	}
	msg, _ := json.Marshal(profile)

	_, _ = fmt.Fprintln(w, string(msg))
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// UpdateEmail
	// UpdateUsername
	// UpdateAvatar
	// whatever else...

	var message models.SuccessOrErrorMessage

	var newProfile models.UserUpdateInfo
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		message.Send(w, http.StatusBadRequest, "JSON parsing error")
		return
	}

	id := userGuid(r)
	user, _ := userStorage.Get(id)
	user.Username = newProfile.Username
	user.Email = newProfile.Email

	message.Send(w, http.StatusOK, "Profile successfully updated")
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := userGuid(r)
	var message models.SuccessOrErrorMessage

	err := userStorage.Delete(id)
	if err != nil {
		message.Send(w, http.StatusInternalServerError, "Error while deleting profile")
	}

	message.Send(w, http.StatusOK, "Profile "+string(id)+" deleted successfully")
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(int64(5 * MB))
	var message models.SuccessOrErrorMessage
	file, _, err := r.FormFile("avatar")
	if err != nil {
		message.Send(w, http.StatusBadRequest, "Error reading file")
		return
	}
	defer file.Close()

	data := bytes.NewBuffer(nil)
	_, err = io.Copy(data, file)
	if err != nil {
		panic(err)
	}

	id := userGuid(r)
	user, _ := userStorage.Get(id)

	err = ioutil.WriteFile(`/home/daniknik/colors_static/`+user.Username+`.png`, data.Bytes(), 0666)
	if err != nil {
		panic(err)
	}

	user.AvatarUrl = "/img/" + user.Username + `.png`

	message.Send(w, http.StatusOK, user.AvatarUrl)
}

//сервисный метод чтобы понимать что творится в хранилище юзеров
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	message, _ := json.Marshal(userStorage.Data)
	_, _ = fmt.Fprintln(w, string(message))
}
