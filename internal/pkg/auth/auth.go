package auth

import (
    "crypto/sha256"
    "database/sql"
    "fmt"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
    "net/http"
    "strings"
    "time"
)

func CreateAuthCookie(data models.JwtData, lifetime time.Duration) *http.Cookie {
    jwtStr, err := data.Marshal(lifetime, secret)
    if err != nil {
        return &http.Cookie{}
    }

    return &http.Cookie{
        Name:     CookieName,
        Value:    jwtStr,
        Expires:  time.Now().Add(lifetime),
        HttpOnly: true,
    }
}

func CheckJwt(token string) (models.JwtData, error) {
    data := models.JwtData{}
    err := data.UnMarshal(token, secret)
    return data, err
}

func SignUp(signUpData models.SingUpData) (models.JwtData, error) {
    incorrectFields := signUpData.Check()
    if len(incorrectFields) > 0 {
        return models.JwtData{}, fmt.Errorf("incorrect: " + strings.Join(incorrectFields, ", "))
    }

    passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(signUpData.Password)))

    id, err := db.AuthCreate(db.AuthData{
        Email:    signUpData.Email,
        Username: signUpData.Username,
        PassHash: passHash,
    })
    if err != nil {
        return models.JwtData{}, err
    }

    return models.JwtData{
        Id:       id,
        Email:    signUpData.Email,
        Username: signUpData.Username,
    }, nil
}

func SignIn(signInData models.SignInData) (data models.JwtData, err error) {
    var userData db.AuthData
    passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(signInData.Password)))

    isEmail := models.CheckEmail(signInData.Login)
    if isEmail {
        userData, err = db.AuthFindByEmailAndPassHash(signInData.Login, passHash)
        if err != nil {
            if err == sql.ErrNoRows {
                return data, fmt.Errorf("incorrect login or password")
            }
            return
        }
    }

    isUsername := !isEmail && models.CheckUsername(signInData.Login)
    if isUsername {
        userData, err = db.AuthFindByNicknameAndPassHash(signInData.Login, passHash)
        if err != nil {
            if err == sql.ErrNoRows {
                return data, fmt.Errorf("incorrect login or password")
            }
            return
        }
    }

    if !isEmail && !isUsername {
        return data, fmt.Errorf("incorrect login")
    }

    return models.JwtData{
        Id:       userData.Id,
        Email:    userData.Email,
        Username: userData.Username,
    }, nil
}

func UpdatePassword(id int64, passwordData models.UpdatePasswordData) error {
    incorrectFields := passwordData.Check()
    if len(incorrectFields) > 0 {
        return fmt.Errorf("incorrect: " + strings.Join(incorrectFields, ", "))
    }

    passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(passwordData.NewPassword)))
    return db.AuthUpdatePassword(id, passHash)
}

func UpdateAuth(id int64, userData models.UpdateUserData) (models.JwtData, error) {
    incorrectFields := userData.Check()
    if len(incorrectFields) > 0 {
        return models.JwtData{}, fmt.Errorf("incorrect: " + strings.Join(incorrectFields, ", "))
    }

    err := db.AuthUpdateData(db.AuthData{
        Id:       id,
        Email:    userData.Email,
        Username: userData.Username,
    })
    if err != nil {
        return models.JwtData{}, err
    }

    return models.JwtData{
        Id:       id,
        Email:    userData.Email,
        Username: userData.Username,
    }, nil
}

func RemoveAuth(id int64, removeData models.RemoveUserData) error {
    incorrectFields := removeData.Check()
    if len(incorrectFields) > 0 {
        return fmt.Errorf("incorrect: " + strings.Join(incorrectFields, ", "))
    }

    passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(removeData.Password)))
    return db.AuthRemove(id, passHash)
}
