package controllers

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

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
// @router /avatar [post]
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.SendMessage(w, http.StatusUnauthorized, "not signed in")
		return
	}

	_ = r.ParseMultipartForm(int64(5 * MByte))
	file, header, err := r.FormFile("avatar")
	if err != nil {
		models.SendMessage(w, http.StatusBadRequest, "error reading file: "+err.Error())
		return
	}
	defer file.Close()

	id := jwtData(r).Id
	ext := strings.Split(header.Filename, ".")[1]
	err = fileVault.UploadFile(file, strconv.Itoa(int(id)), ext)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, "impossible save file: "+err.Error())
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
