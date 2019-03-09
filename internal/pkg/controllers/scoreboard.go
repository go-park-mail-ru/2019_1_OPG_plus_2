package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"github.com/gorilla/mux"
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

func ScoreBoardByPage(w http.ResponseWriter, r *http.Request) {

	pathVariables := mux.Vars(r)
	var message models.SuccessOrErrorMessage
	if pathVariables == nil {
		message.Send(w, http.StatusBadRequest, "Bad query")
		return
	}
	pageVar, ok := pathVariables["page"]
	if !ok {
		message.Send(w, http.StatusBadRequest, "Bad query")
		return
	}

	page, _ := strconv.ParseInt(pageVar, 10, 32)

	lbound := int(page)*pageSize - 1
	for lbound > len(scoreboard) {
		lbound -= pageSize
		if lbound < 0 {
			lbound = 0
		}
	}

	rbound := int(page) * pageSize
	if rbound > len(scoreboard) {
		rbound = len(scoreboard)
	}

	scoreboardPage := scoreboard[lbound:rbound]
	msg, err := json.Marshal(scoreboardPage)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = fmt.Fprint(w, string(msg))
}
