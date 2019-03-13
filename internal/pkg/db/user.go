package db

import (
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

func GetUser(id int64) (userData models.UserData, err error) {
    row, err := QueryRow("SELECT a.id, a.username, a.email, c.score, c.games, c.win, c.lose FROM "+
        authDbName+"."+authUsersTable+" AS a JOIN "+coreDbName+"."+coreUsersTable+" AS c ON a.id = c.id WHERE id = ?", id)
    if err != nil {
        return
    }
    err = row.Scan(&userData.Id, &userData.Username, &userData.Email, &userData.Score, &userData.Games, &userData.Win, &userData.Lose)
    return
}

func GetScoreboard(limit, offset int64) (usersData []models.ScoreboardUserData, err error) {
    rows, err := Query("SELECT a.id, a.username, c.score FROM "+authDbName+"."+authUsersTable+" AS a JOIN "+
        coreDbName+"."+coreUsersTable+" AS c ON a.id = c.id ORDER BY c.score DESC LIMIT ? OFFSET ?", limit, offset)
    if err != nil {
        return
    }

    for rows.Next() {
        userData := models.ScoreboardUserData{}
        err = rows.Scan(&userData.Id, &userData.Username, &userData.Score)
        if err != nil {
            return
        }
        usersData = append(usersData, userData)
    }
    return
}
