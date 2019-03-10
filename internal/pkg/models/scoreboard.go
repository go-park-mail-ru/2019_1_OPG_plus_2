package models

type ScoreboardRecord struct {
	Position int    `json:"position, number" example:"1"`
	Username string `json:"username, string" example:"XxX__NaGiBaToR__XxX"`
	Score    int    `json:"score, number" example:"314159"`
}
