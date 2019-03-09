package models

type ScoreboardRecord struct {
	Position int    `json:"position, number"`
	Username string `json:"username, string"`
	Score    int    `json:"score, number"`
}
