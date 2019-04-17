package db

import (
	"database/sql"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

type ProfileData struct {
	Id     int64  `json:"id"`
	Avatar string `json:"avatar, string" example:"<some avatar url>"`
	Score  int64  `json:"score, number"`
	Games  int64  `json:"games, number"`
	Win    int64  `json:"win, number"`
	Lose   int64  `json:"lose, number"`
}

func ProfileCreate(data ProfileData) (err error) {
	id, err := isExists(CoreDbName, CoreUsersTable, "id = ?", data.Id)
	if err != nil {
		return
	}
	if id != 0 {
		return models.AlreadyExists
	}

	_, err = insert(CoreDbName, CoreUsersTable, "id, avatar", "?, ?", data.Id, data.Avatar)
	return
}

func ProfileFindById(id int64) (data ProfileData, err error) {
	row, err := findRowBy(CoreDbName, CoreUsersTable, "avatar, score, games, win, lose", "id = ?", id)
	if err != nil {
		return
	}
	data.Id = id
	err = row.Scan(&data.Avatar, &data.Score, &data.Games, &data.Win, &data.Lose)
	if err == sql.ErrNoRows {
		return data, models.NotFound
	}
	return
}

func ProfileUpdateData(data ProfileData) error {
	id, err := isExists(CoreDbName, CoreUsersTable, "id = ?", data.Id)
	if err != nil {
		return err
	}
	if id == 0 {
		return models.NotFound
	}

	_, err = updateBy(CoreDbName, CoreUsersTable, "score = ?, games = ?, win = ?, lose = ?", "id = ?",
		data.Score, data.Games, data.Win, data.Lose, data.Id)
	return err
}

func ProfileUpdateAvatar(id int64, avatar string) error {
	id, err := isExists(CoreDbName, CoreUsersTable, "id = ?", id)
	if err != nil {
		return err
	}
	if id == 0 {
		return models.NotFound
	}

	_, err = updateBy(CoreDbName, CoreUsersTable, "avatar = ?", "id = ?", avatar, id)
	return err
}

func ProfileRemove(id int64) error {
	id, err := isExists(CoreDbName, CoreUsersTable, "id = ?", id)
	if err != nil {
		return err
	}
	if id == 0 {
		return models.NotFound
	}

	_, err = removeBy(CoreDbName, CoreUsersTable, "id = ?", id)
	return err
}

func ProfileTruncate() error {
	return truncate(CoreDbName, CoreUsersTable)
}
