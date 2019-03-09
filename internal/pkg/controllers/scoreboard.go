package controllers

import (
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
	"net/http"
)

var scoreboard = []models.ScoreboardRecord{
	models.ScoreboardRecord{
		Username: "user1",
		Position: 1,
		Score:    1000,
	},
	models.ScoreboardRecord{
		Username: "user2",
		Position: 2,
		Score:    500,
	},
}

func ScoreBoardByPage(w http.ResponseWriter, r *http.Request) {

}
