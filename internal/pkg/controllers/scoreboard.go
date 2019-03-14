package controllers

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"strconv"
)

// GetScoreboard godoc
// @title Get scoreboard page
// @summary Produces scoreboard page with {limit} and {offset}
// @description This method provides client with scoreboard limited with {limit} entries per page and offset of {offset} from the first position
// @tags scoreboard
// @produce json
// @param limit query int false "Entries per page"
// @param offset query int false "Entries from the first position"
// @success 200 {array} models.ScoreboardUserData
// @router /users [get]
func GetScoreboard(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit < 1 {
		limit = pageSize
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	usersData, err := db.GetScoreboard(limit, (page-1)*limit)
	if err != nil {
		models.SendMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	models.SendMessageWithData(w, http.StatusOK, "users found", usersData)
}
