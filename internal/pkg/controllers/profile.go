package controllers

import (
    "encoding/json"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/auth"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/profile"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/util/fileStorage"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var MB = 1 << 20

var fileVault = fileStorage.NewLocalFileStorage("/home/daniknik/colors_static")

// CreateProfile godoc
// @title Create profile
// @summary Registers user
// @description This method creates records about new user in auth-bd and profile-db and then sends cookie to user in order to identify
// @tags profile
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User profile data"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /profile [post]
func CreateProfile(w http.ResponseWriter, r *http.Request) {
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

    jwtData, err := profile.CreateProfile(signUpData)
    if err != nil {
        models.SendMessage(w, http.StatusUnauthorized, err.Error())
        return
    }

    http.SetCookie(w, auth.CreateAuthCookie(jwtData, 30*24*time.Hour))
    models.SendMessage(w, http.StatusOK, "signed up")
}

// GetProfile godoc
// @title Get profile
// @summary Produces user profile info
// @description This method provides client with profile data, matching required ID
// @tags profile
// @accept json
// @produce json
// @param id path int true "Profile ID"
// @success 200 {object} models.ProfileData
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 404 {object} models.SuccessOrErrorMessage
// @router /profile/{id} [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
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

    userData, err := profile.GetProfile(id)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }
    msg, err := json.Marshal(userData)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessage(w, http.StatusOK, string(msg))
}

// UpdateProfile godoc
// @title Update profile
// @summary Updates client's profile
// @description This method updates info in profile and auth-db record of user, who is making a query
// @tags profile
// @accept json
// @produce json
// @param profile_data body models.SingUpData true "User new profile data"
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 401 {object} models.SuccessOrErrorMessage
// @router /profile [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
    if !isAuth(r) {
        models.SendMessage(w, http.StatusUnauthorized, "not signed in")
        return
    }

    var newProfile models.UserData
    err := json.NewDecoder(r.Body).Decode(&newProfile)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, "incorrect JSON")
        return
    }

    newProfile.Id = jwtData(r).Id
    err = profile.UpdateProfile(newProfile)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessage(w, http.StatusOK, "profile updated")
}

// DeleteProfile godoc
// @title Delete profile
// @summary Deletes profile and user of client
// @description This method deletes all information about user, making a query, including profile, game stats and authorization info
// @tags profile
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /profile [delete]
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
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

    err = profile.RemoveProfile(jwtData(r).Id, removeData)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessage(w, http.StatusOK, "profile removed")
}

// TODO: implement FileStorage interface to inline file operations

// UploadAvatar godoc
// @title Upload new avatar
// @summary Saves new avatar image of client's user
// @description This method saves avatar image in server storage and sets it as clients user avatar
// @tags profile
// @accept png
// @accept jpeg
// @produce json
// @success 200 {object} models.SuccessOrErrorMessage
// @failure 400 {object} models.SuccessOrErrorMessage
// @failure 500 {object} models.SuccessOrErrorMessage
// @router /upload_avatar [post]
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
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
    err = profile.UpdateAvatar(id, newAvatar)
    if err != nil {
        models.SendMessage(w, http.StatusInternalServerError, err.Error())
        return
    }

    models.SendMessage(w, http.StatusOK, newAvatar)
}
