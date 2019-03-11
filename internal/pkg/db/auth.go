package db

import (
    "database/sql"
    "fmt"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
    _ "github.com/go-sql-driver/mysql"
)

var authDb *sql.DB

func AuthInit() (err error) {
    if authDb != nil {
        return fmt.Errorf("db already initialized")
    }
    authDb, err = sql.Open("mysql", authUsername + ":" + authPassword + "@/" +authDbName)
    return
}

func AuthClose() error {
    return authDb.Close()
}

func AuthQuery(query string, args ...interface{}) (*sql.Rows, error) {
    return authDb.Query(query, args...)
}

func AuthExec(query string, args ...interface{}) (sql.Result, error) {
    return authDb.Exec(query, args...)
}

func AuthCreate(data models.DbUserData) (int64, error) {
    return insert(authDb, authDbName, authUsersTable,"username, email, pass_hash", "?, ?, ?",
        data.Username, data.Email, data.PassHash)
}

func AuthFindById(id int64) (data models.DbUserData, err error) {
    row, err := findRowBy(authDb, authDbName, authUsersTable, "id, username, email, pass_hash", "id = ?", id)
    if err != nil {
        return
    }
    err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
    return
}

func AuthFindByEmailAndPassHash(email string, passHash string) (data models.DbUserData, err error) {
    row, err := findRowBy(authDb, authDbName, authUsersTable, "id, username, email, pass_hash", "email = ? AND pass_hash = ?", email, passHash)
    if err != nil {
        return
    }
    err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
    return
}

func AuthFindByNicknameAndPassHash(username string, passHash string) (data models.DbUserData, err error) {
    row, err := findRowBy(authDb, authDbName, authUsersTable, "id, username, email, pass_hash", "username = ? AND pass_hash = ?", username, passHash)
    if err != nil {
        return
    }
    err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
    return
}

func AuthUpdateData(data models.DbUserData) error {
    _, err := updateBy(authDb, authDbName, authUsersTable,"username = ?, email = ?", "id = ?",
        data.Username, data.Email, data.Id)
    return err
}

func AuthUpdatePassword(id int64, passHash string) error {
    _, err := updateBy(authDb, authDbName, authUsersTable,"pass_hash = ?", "id = ?", passHash, id)
    return err
}

func AuthRemove(id int64, passHash string) error {
    _, err := removeBy(authDb, authDbName, authUsersTable,"id = ? AND pass_hash = ?", id, passHash)
    return err
}
