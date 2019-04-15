package controllers

import (
	"net/http"
	"strconv"

	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

const pageSize = 10

// GetScoreboard godoc
// @title Get scoreboard page
// @summary Produces scoreboard page with {limit} and {offset}
// @description This method provides client with scoreboard limited with {limit} entries per page and offset of {offset} from the first position
// @tags scoreboard
// @produce json
// @param limit query int false "Entries per page"
// @param page query int false "Number of page"
// @success 200 {array} models.ScoreboardUserData
// @failure 500 {object} models.MessageAnswer
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

	users, count, err := db.GetScoreboard(limit, (page-1)*limit)
	if err != nil {
		models.Send(w, http.StatusInternalServerError, models.GetDeveloperErrorAnswer(err.Error()))
		return
	}

	models.Send(w, http.StatusOK, models.GetScoreboardAnswer(models.ScoreboardData{
		Users: users,
		Count: count,
	}))
}
