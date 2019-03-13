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
    //OldPassword     string `json:"old_password" example:"SecretPass1!"`
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

func (data RemoveUserData) Check() (incorrectFields []string) {
    return
}

/********************
 *  IN-OUT MODELS   *
 ********************/

/* USER DATA */

type UserData struct {
    Id       int64  `json:"id, string" example:"1"`
    Username string `json:"username, string" example:"user_test"`
    Email    string `json:"email, string" example:"user_test@test.com"`
    Avatar   string `json:"avatar, string" example:"<some avatar url>"`
    Score    int64  `json:"score, number"`
    Games    int64  `json:"games, number"`
    Win      int64  `json:"win, number"`
    Lose     int64  `json:"lose, number"`
}

/* SCOREBOARD DATA */

type ScoreboardUserData struct {
    Id       int64  `json:"id, string" example:"1"`
    Username string `json:"username, string" example:"XxX__NaGiBaToR__XxX"`
    Score    int64  `json:"score, number" example:"314159"`
}
