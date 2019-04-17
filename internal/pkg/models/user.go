package models

import (
	"regexp"
)

/********************
 *    IN MODELS     *
 ********************/

/* COMMON CHECKERS */

func CheckEmail(email string) bool {
	matched, _ := regexp.MatchString(`^[\w\-.]+@[\w\-.]+\.[a-z]{2,6}$`, email)
	return matched
}

func CheckUsername(username string) bool {
	matched, _ := regexp.MatchString(`^\w+$`, username)
	return matched
}

/* SIGN IN DATA */

type SignInData struct {
	Login    string `json:"login" example:"test@mail.ru"`
	Password string `json:"password" example:"Qwerty123"`
}

func (data SignInData) Check() (incorrectFields []string) {
	if !CheckEmail(data.Login) && !CheckUsername(data.Login) {
		incorrectFields = append(incorrectFields, "login")
	}
	if data.Password == "" {
		incorrectFields = append(incorrectFields, "password")
	}
	return
}

/* SIGN UP DATA */

type SignUpData struct {
	Email    string `json:"email" example:"user_test@test.com"`
	Username string `json:"username" example:"user_test"`
	Password string `json:"password" example:"SecretPass1!"`
}

func (data SignUpData) Check() (incorrectFields []string) {
	if !CheckEmail(data.Email) {
		incorrectFields = append(incorrectFields, "email")
	}
	if !CheckUsername(data.Username) {
		incorrectFields = append(incorrectFields, "username")
	}
	if data.Password == "" {
		incorrectFields = append(incorrectFields, "password")
	}
	return
}

/* UPDATE USER DATA */

type UpdateUserData struct {
	Email    string `json:"email" example:"user_test@test.com"`
	Username string `json:"username" example:"user_test"`
}

func (data UpdateUserData) Check() (incorrectFields []string) {
	if !CheckEmail(data.Email) {
		incorrectFields = append(incorrectFields, "email")
	}
	if !CheckUsername(data.Username) {
		incorrectFields = append(incorrectFields, "username")
	}
	return
}

/* UPDATE PASSWORD DATA */

type UpdatePasswordData struct {
	NewPassword     string `json:"new_password" example:"SecretPass2!"`
	PasswordConfirm string `json:"password_confirm" example:"SecretPass2!"`
}

func (data UpdatePasswordData) Check() (incorrectFields []string) {
	if data.NewPassword == "" {
		incorrectFields = append(incorrectFields, "new_password")
	}
	if data.PasswordConfirm != data.NewPassword {
		incorrectFields = append(incorrectFields, "password_confirm")
	}
	return
}

/* REMOVE USER DATA */

type RemoveUserData struct {
	Password string `json:"password" example:"SecretPass1!"`
}

func (data RemoveUserData) Check() (incorrectFields []string) {
	if data.Password == "" {
		incorrectFields = append(incorrectFields, "password")
	}
	return
}

/********************
 *    OUT MODELS    *
 ********************/

/* USER DATA */

type UserData struct {
	Id       int64  `json:"id, string" example:"1"`
	Username string `json:"username, string" example:"user_test"`
	Email    string `json:"email, string" example:"user_test@test.com"`
	Avatar   string `json:"avatar, string" example:"<some avatar url>"`

	Score int64 `json:"score, number"`
	Games int64 `json:"games, number"`
	Win   int64 `json:"win, number"`
	Lose  int64 `json:"lose, number"`
}

/* SCOREBOARD DATA */

type ScoreboardUserData struct {
	Id       int64  `json:"id, string" example:"1"`
	Username string `json:"username, string" example:"XxX__NaGiBaToR__XxX"`
	Avatar   string `json:"avatar, string" example:"<some avatar url>"`
	Score    int64  `json:"score, number" example:"314159"`
}

type ScoreboardData struct {
	Users []ScoreboardUserData `json:"users"`
	Count uint64               `json:"count" example:"123"`
}
