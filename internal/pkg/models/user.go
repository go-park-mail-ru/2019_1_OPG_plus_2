package models

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	if !CheckEmail(data.Login) {
		return []string{"login"}
	}
	if !CheckUsername(data.Login) {
		return []string{"login"}
	}
	return
}

/* SIGN UP DATA */

type SingUpData struct {
	Email    string `json:"email" example:"user_test@test.com"`
	Username string `json:"username" example:"user_test"`
	Password string `json:"password" example:"SecretPass1!"`
	Avatar   string `json:"avatar, string" example:"<some avatar url>"`
}

func (data SingUpData) Check() (incorrectFields []string) {
	if !CheckEmail(data.Email) {
		incorrectFields = append(incorrectFields, "email")
	}
	if !CheckUsername(data.Username) {
		incorrectFields = append(incorrectFields, "username")
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
	if data.PasswordConfirm != data.NewPassword {
		incorrectFields = append(incorrectFields, "password_confirm")
	}
	return
}

/* REMOVE USER DATA */

type RemoveUserData struct {
	Password string `json:"password" example:"SecretPass1!"`
}

//  TODO: валидация приходящих данных на предмет поля password
func (data RemoveUserData) Check() (incorrectFields []string) {
	return
}

/********************
 *    OUT MODELS    *
 ********************/

/* SIGN IN */

type SignInAnswer struct {
	Login string `json:"login" example:"test@mail.ru"`
}

type SignInAnswerMessage struct {
	AnswerMessage
	Data SignInAnswer `json:"data"`
}

func (message SignInAnswerMessage) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendSignInAnswer(w http.ResponseWriter, status int, message string, data SignInAnswer) {
	SignInAnswerMessage{
		AnswerMessage: AnswerMessage{
			Status:  status,
			Message: message,
		},
		Data: data,
	}.Send(w)
}

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

type UserDataAnswerMessage struct {
	AnswerMessage
	Data UserData `json:"data"`
}

func (message UserDataAnswerMessage) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendUserDataAnswer(w http.ResponseWriter, status int, message string, data UserData) {
	UserDataAnswerMessage{
		AnswerMessage: AnswerMessage{
			Status:  status,
			Message: message,
		},
		Data: data,
	}.Send(w)
}

/* SCOREBOARD DATA */

type ScoreboardUserData struct {
	Id       int64  `json:"id, string" example:"1"`
	Username string `json:"username, string" example:"XxX__NaGiBaToR__XxX"`
	Score    int64  `json:"score, number" example:"314159"`
}

type ScoreboardAnswerMessage struct {
	AnswerMessage
	Data []ScoreboardUserData `json:"data"`
}

func (message ScoreboardAnswerMessage) Send(w http.ResponseWriter) {
	w.WriteHeader(message.Status)
	msg, _ := json.Marshal(message)
	_, _ = fmt.Fprintln(w, string(msg))
}

func SendScoreboardAnswer(w http.ResponseWriter, status int, message string, data []ScoreboardUserData) {
	ScoreboardAnswerMessage{
		AnswerMessage: AnswerMessage{
			Status:  status,
			Message: message,
		},
		Data: data,
	}.Send(w)
}
