package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"2019_1_OPG_plus_2/internal/pkg/models"
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
		return id, models.AlreadyExists
	}

	return insert(authDbName, authUsersTable, "username, email, pass_hash", "?, ?, ?",
		data.Username, data.Email, data.PassHash)
}

// For future use
//
// func AuthFindByUsername(username string) (data AuthData, err error) {
// 	row, err := findRowBy(authDbName, authUsersTable, "id, username, email, pass_hash", "username = ?", username)
// 	if err != nil {
// 		return
// 	}
// 	err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
// 	return
// }

func AuthFindByEmailAndPassHash(email string, passHash string) (data AuthData, err error) {
	row, err := findRowBy(authDbName, authUsersTable, "id, username, email, pass_hash", "email = ? AND pass_hash = ?", email, passHash)
	if err != nil {
		return
	}
	err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
	if err == sql.ErrNoRows {
		return data, models.NotFound
	}
	return
}

func AuthFindByUsernameAndPassHash(username string, passHash string) (data AuthData, err error) {
	row, err := findRowBy(authDbName, authUsersTable, "id, username, email, pass_hash", "username = ? AND pass_hash = ?", username, passHash)
	if err != nil {
		return
	}
	err = row.Scan(&data.Id, &data.Username, &data.Email, &data.PassHash)
	if err == sql.ErrNoRows {
		return data, models.NotFound
	}
	return
}

func AuthUpdateData(data AuthData) error {
	id, err := isExists(authDbName, authUsersTable, "id = ?", data.Id)
	if err != nil {
		return err
	}
	if id == 0 {
		return models.NotFound
	}

	_, err = updateBy(authDbName, authUsersTable, "username = ?, email = ?", "id = ?", data.Username, data.Email, data.Id)
	if mysqlError, ok := err.(*mysql.MySQLError); ok {
		if mysqlError.Number == 1062 {
			return models.AlreadyExists
		}
	}
	return err
}

func AuthUpdatePassword(id int64, passHash string) error {
	id, err := isExists(authDbName, authUsersTable, "id = ?", id)
	if err != nil {
		return err
	}
	if id == 0 {
		return models.NotFound
	}

	_, err = updateBy(authDbName, authUsersTable, "pass_hash = ?", "id = ?", passHash, id)
	return err
}

func AuthRemove(id int64, passHash string) error {
	id, err := isExists(authDbName, authUsersTable, "id = ?", id)
	if err != nil {
		return err
	}
	if id == 0 {
		return models.NotFound
	}

	count, err := removeBy(authDbName, authUsersTable, "id = ? AND pass_hash = ?", id, passHash)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect password")
	}
	return nil
}

func AuthTruncate() error {
	return truncate(authDbName, authUsersTable)
}
