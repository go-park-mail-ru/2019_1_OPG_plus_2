package db

import (
    "fmt"
)

type ProfileData struct {
    Id     int64  `json:"id"`
    Avatar string `json:"avatar, string" example:"<some avatar url>"`
    Score  int64    `json:"score, number"`
    Games  int64    `json:"games, number"`
    Win    int64    `json:"win, number"`
    Lose   int64    `json:"lose, number"`
}

func ProfileCreate(data ProfileData) (err error) {
    id, err := isExists(coreDbName, coreUsersTable, "id = ?", data.Id)
    if err != nil {
        return
    }
    if id != 0 {
        return fmt.Errorf("user already exists")
    }

    id, err = insert(coreDbName, coreUsersTable,"id, avatar", "?, ?", data.Id, data.Avatar)
    return
}

func ProfileFindById(id int64) (data ProfileData, err error) {
    row, err := findRowBy(coreDbName, coreUsersTable, "avatar, score, games, win, lose", "id = ?", id)
    if err != nil {
        return
    }
    data.Id = id
    err = row.Scan(&data.Avatar, &data.Score, &data.Games, &data.Win, &data.Lose)
    return
}

func ProfileUpdateData(data ProfileData) error {
    count, err := updateBy(coreDbName, coreUsersTable,"score = ?, games = ?, win = ?, lose = ?", "id = ?",
        data.Score, data.Games, data.Win, data.Lose, data.Id)
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}

func ProfileUpdateAvatar(id int64, avatar string) error {
    count, err := updateBy(coreDbName, coreUsersTable,"avatar = ?", "id = ?", avatar, id)
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}

func ProfileRemove(id int64) error {
    count, err := removeBy(coreDbName, coreUsersTable,"id = ?", id)
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}
