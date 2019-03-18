package controllers

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/util/fileStorage"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"
)

const mByte = 1 << 20

var StaticPath, _ = filepath.Abs("./static")
var fileVault = fileStorage.NewLocalFileStorage(StaticPath)

func isImage(header textproto.MIMEHeader) bool {
	return strings.HasPrefix(header.Get("Content-Type"), "image/")
}

// UploadAvatar godoc
// @title Upload new avatar
// @summary Saves new avatar image of client's user
// @description This method saves avatar image in server storage and sets it as clients user avatar
// @tags user
// @accept png
// @accept jpeg
// @produce json
// @success 200 {object} models.MessageAnswer
// @failure 400 {object} models.MessageAnswer
// @failure 500 {object} models.MessageAnswer
// @router /avatar [post]
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if !isAuth(r) {
		models.Send(w, http.StatusUnauthorized, models.NotSignedInAnswer)
		return
	}

	err := r.ParseMultipartForm(int64(5 * mByte))
	if err != nil {
		models.Send(w, http.StatusBadRequest, models.FileTooBigAnswer)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		models.Send(w, http.StatusBadRequest, models.ImpossibleReadFileAnswer)
		return
	}
	defer file.Close()

	if !isImage(header.Header) {
		models.Send(w, http.StatusBadRequest, models.NotImageAnswer)
		return
	}

	nameParts := strings.Split(header.Filename, ".")
	ext := nameParts[len(nameParts)-1]

	newName := strconv.FormatInt(jwtData(r).Id, 10)
	err = fileVault.UploadFile(file, newName, ext)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.ImpossibleSaveFileAnswer)
		return
	}

	url := "/static/" + newName + "." + ext
	err = db.ProfileUpdateAvatar(jwtData(r).Id, url)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.MessageAnswer{Status: 500, Message: err.Error()})
		return
	}

	models.Send(w, http.StatusOK, models.MessageAnswer{Status: 100, Message: url})
}
