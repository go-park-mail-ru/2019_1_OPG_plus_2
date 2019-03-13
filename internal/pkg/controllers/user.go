package controllers

import (
    "encoding/json"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/user"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/util/fileStorage"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var MB = 1 << 20

var fileVault = fileStorage.NewLocalFileStorage("/home/daniknik/colors_static")

// CreateUser godoc
// @title Create user
// @summary Registers user
// @description This method creates records about new user in auth-bd and user-db and then sends cookie to user in order to identify
// @tags user
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User user data"
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessage
// @failure 401 {object} models.AnswerMessage
// @router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
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

    jwtData, err := user.CreateUser(signUpData)
    if err != nil {
        models.SendMessage(w, http.StatusUnauthorized, err.Error())
        return
    }

    http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
    models.SendMessage(w, http.StatusOK, "signed up")
}

// GetUser godoc
// @title Get user
// @summary Produces user user info
// @description This method provides client with user data, matching required ID
// @tags user
// @accept json
// @produce json
// @param id path int true "Profile ID"
// @success 200 {object} models.ProfileData
// @failure 400 {object} models.AnswerMessage
// @failure 404 {object} models.AnswerMessage
// @router /user/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
    var id int64
    var err error

    pathVariables := mux.Vars(r)
    if pathVariables == nil {
        id = jwtData(r).Id
    } else {
        pathId, ok := pathVariables["id"]
        if !ok {
            models.SendMessage(w, http.StatusBadRequest, "no id in query")
            return
        }
        id, err = strconv.ParseInt(pathId, 10, 64)
        if !ok {
            models.SendMessage(w, http.StatusBadRequest, "incorrect id in query")
            return
        }
    }

    userData, err := user.GetUser(id)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessageWithData(w, http.StatusOK, "user found", userData)
}

// UpdateAuth godoc
// @title Update user
// @summary Updates client's user
// @description This method updates info in user and auth-db record of user, who is making a query
// @tags user
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User new user data"
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessage
// @failure 401 {object} models.AnswerMessage
// @router /user [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
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

    jwtData, err := user.UpdateUser(jwtData(r).Id, updateData)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
    models.SendMessage(w, http.StatusOK, "user updated")
}

// RemoveAuth godoc
// @title Delete user
// @summary Deletes user and user of client
// @description This method deletes all information about user, making a query, including user, game stats and authorization info
// @tags user
// @produce json
// @success 200 {object} models.AnswerMessage
// @failure 500 {object} models.AnswerMessage
// @router /user [delete]
func RemoveUser(w http.ResponseWriter, r *http.Request) {
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

    err = user.RemoveUser(jwtData(r).Id, removeData)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessage(w, http.StatusOK, "user removed")
}

// TODO: implement FileStorage interface to inline file operations

// UploadAvatar godoc
// @title Upload new avatar
// @summary Saves new avatar image of client's user
// @description This method saves avatar image in server storage and sets it as clients user avatar
// @tags user
// @accept png
// @accept jpeg
// @produce json
// @success 200 {object} models.AnswerMessage
// @failure 400 {object} models.AnswerMessage
// @failure 500 {object} models.AnswerMessage
// @router /upload_avatar [post]
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
    if !isAuth(r) {
        models.SendMessage(w, http.StatusUnauthorized, "not signed in")
        return
    }

    _ = r.ParseMultipartForm(int64(5 * MB))
    file, header, err := r.FormFile("avatar")
    if err != nil {
        models.SendMessage(w, http.StatusBadRequest, "error reading file")
        return
    }
    defer file.Close()

    id := jwtData(r).Id
    ext := strings.Split(header.Filename, ".")[1]
    err = fileVault.UploadFile(file, strconv.Itoa(int(id)), ext)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, "impossible save file")
        return
    }

    newAvatar := "/static/" + strconv.Itoa(int(id)) + "." + ext
    err = db.ProfileUpdateAvatar(id, newAvatar)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessage(w, http.StatusOK, newAvatar)
}
