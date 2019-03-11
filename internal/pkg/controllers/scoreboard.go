package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
	"strconv"
)

var scoreboard = []models.ScoreboardRecord{
	{
		Username: "user1",
		Position: 1,
		Score:    1000,
	},
	{
		Username: "user2",
		Position: 2,
		Score:    500,
	},
}

var pageSize = 10

// ScoreBoardByPage godoc
// @title Get scoreboard page
// @summary Produces scoreboard page with {limit} and {offset}
// @description This method provides client with scoreboard limited with {limit} entries per page and offset of {offset} from the first position
// @tags scoreboard
// @produce json
// @param limit query int false "Entries per page"
// @param offset query int false "Entries from the first position"
// @success 200 {array} models.ScoreboardRecord
// @router /profiles/score [get]
func ScoreBoardByPage(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = pageSize
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	upperBound := offset + limit
	if upperBound > len(scoreboard) {
		upperBound = len(scoreboard)
	}
	scoreboardPage := scoreboard[offset:upperBound]
	msg, err := json.Marshal(scoreboardPage)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = fmt.Fprint(w, string(msg))
}
