package db

import (
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type AuthData struct {
    Id       int64  `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    PassHash string `json:"pass_hash"`
}

func AuthCreate(data AuthData) (id int64, err error) {
    id, err = isExists(authDbName, authUsersTable, "email = ? OR username = ?", data.Email, data.Username)
    if err != nil {
        return
    }
    if id != 0 {
        return id, fmt.Errorf("user already exists")
    }

    return insert(authDbName, authUsersTable,"username, email, pass_hash", "?, ?, ?",
        data.Username, data.Email, data.PassHash)
}

func AuthFindById(id int64) (data AuthData, err error) {
    row, err := findRowBy(authDbName, authUsersTable, "id, username, email, pass_hash", "id = ?", id)
    if err != nil {
        return
    }
    err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
    return
}

func AuthFindByEmailAndPassHash(email string, passHash string) (data AuthData, err error) {
    row, err := findRowBy(authDbName, authUsersTable, "id, username, email, pass_hash", "email = ? AND pass_hash = ?", email, passHash)
    if err != nil {
        return
    }
    err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
    return
}

func AuthFindByNicknameAndPassHash(username string, passHash string) (data AuthData, err error) {
    row, err := findRowBy(authDbName, authUsersTable, "id, username, email, pass_hash", "username = ? AND pass_hash = ?", username, passHash)
    if err != nil {
        return
    }
    err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
    return
}

func AuthUpdateData(data AuthData) error {
    count, err := updateBy(authDbName, authUsersTable,"username = ?, email = ?", "id = ?",
        data.Username, data.Email, data.Id)
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}

func AuthUpdatePassword(id int64, passHash string) error {
    count, err := updateBy(authDbName, authUsersTable,"pass_hash = ?", "id = ?", passHash, id)
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}

func AuthRemove(id int64, passHash string) error {
    count, err := removeBy(authDbName, authUsersTable,"id = ? AND pass_hash = ?", id, passHash)
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("user not found")
    }
    return nil
}
