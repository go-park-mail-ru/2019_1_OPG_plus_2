package db

import (
	"database/sql"

	"2019_1_OPG_plus_2/internal/pkg/models"
)

func GetUser(id int64) (userData models.UserData, err error) {
	row, err := QueryRow("SELECT a.id, a.username, a.email, c.avatar, c.score, c.games, c.win, c.lose FROM "+
		authDbName+"."+authUsersTable+" AS a JOIN "+coreDbName+"."+coreUsersTable+" AS c ON a.id = c.id WHERE a.id = ?", id)
	if err != nil {
		return
	}

	err = row.Scan(&userData.Id, &userData.Username, &userData.Email, &userData.Avatar, &userData.Score,
		&userData.Games, &userData.Win, &userData.Lose)
	if err == sql.ErrNoRows {
		return userData, models.NotFound
	}
	return
}

func GetScoreboard(limit, offset int64) (usersData []models.ScoreboardUserData, count uint64, err error) {
	row, err := QueryRow("SELECT COUNT(id) FROM " + authDbName + "." + authUsersTable)
	if err != nil {
		return
	}
	err = row.Scan(&count)

	rows, err := Query("SELECT a.id, a.username, c.avatar, c.score FROM "+authDbName+"."+authUsersTable+" AS a JOIN "+
		coreDbName+"."+coreUsersTable+" AS c ON a.id = c.id ORDER BY c.score DESC, c.win DESC, c.id LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		userData := models.ScoreboardUserData{}
		err = rows.Scan(&userData.Id, &userData.Username, &userData.Avatar, &userData.Score)
		if err != nil {
			return
		}
		usersData = append(usersData, userData)
	}
	return
}
