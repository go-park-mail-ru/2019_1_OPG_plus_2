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
		ID:        1,
		Score:     228,
		AvatarUrl: "<user1_avatar_url>",
	})

	_ = userStorage.Set("2", &models.UserProfile{
		ID:        2,
		Score:     1337,
		AvatarUrl: "<user2_avatar_url>",
	})
}

//TODO: профайл поменять в соответствии с фоткой: оставить id(number), score, avatar, статы остальное лежит в jwt
//TODO: CreateProfile вызывает controllers/auth.go 66:72 строки, они создают пользователя, я из jwt беру данные (id) и создаю профиль
//TODO: в ScoreBoard пагинация по limit offset

type ProfileData struct {
	Username string `json:"username, string"`
	Email    string `json:"email, string"`
	models.UserProfile
}

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.UserProfile
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		models.SendMessage(w, http.StatusBadRequest, "JSON parsing error")
		return
	}

	id := jwtData(r).Id
	profile.ID = id
	_ = userStorage.Set(string(id), &profile)
	models.SendMessage(w, http.StatusOK, "Profile successfully created")
}

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

	profile, err := userStorage.Get(id)
	if err != nil {
		models.SendMessage(w, http.StatusNotFound, "User not found")
		return
	}
	msg, _ := json.Marshal(profile)

	_, _ = fmt.Fprintln(w, string(msg))
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// UpdateEmail
	// UpdateUsername
	// whatever else...

	var newProfile models.UserUpdateInfo
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		models.SendMessage(w, http.StatusBadRequest, "JSON parsing error")
		return
	}

	//id := strconv.Itoa(jwtData(r).Id)
	//user, _ := userStorage.Get(id)
	//user.Username = newProfile.Username
	//user.Email = newProfile.Email

	//auth.UpdateUser(newProfile)

	models.SendMessage(w, http.StatusOK, "Profile successfully updated")
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := strconv.Itoa(jwtData(r).Id)

	err := userStorage.Delete(id)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "Error while deleting profile")
	}

	models.SendMessage(w, http.StatusOK, "Profile "+string(id)+" deleted successfully")
}

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

	id := strconv.Itoa(jwtData(r).Id)
	fmt.Println(id)
	user, _ := userStorage.Get(id)

	err = ioutil.WriteFile(`/home/daniknik/colors_static/`+id+`.png`, data.Bytes(), 0666)
	if err != nil {
		panic(err)
	}

	user.AvatarUrl = "/img/" + id + `.png`

	models.SendMessage(w, http.StatusOK, user.AvatarUrl)
}

//сервисный метод чтобы понимать что творится в хранилище юзеров
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	message, _ := json.Marshal(userStorage.Data)
	_, _ = fmt.Fprintln(w, string(message))
}
