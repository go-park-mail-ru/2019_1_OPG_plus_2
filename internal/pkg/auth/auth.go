package auth

import (
    "crypto/sha256"
    "database/sql"
    "fmt"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/db"
    "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
    "net/http"
    "regexp"
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

func CreateUser(userData models.UserData) (models.JwtData, error) {
    if matched, _ := regexp.MatchString(`^[\w\-.]+@[\w\-.]+\.[a-z]{2,6}$`, userData.Email); !matched {
        return models.JwtData{}, fmt.Errorf("incorrect email")
    }
    if matched, _ := regexp.MatchString(`^\w+$`, userData.Username); !matched {
        return models.JwtData{}, fmt.Errorf("incorrect nickname")
    }
    passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(userData.Password)))

    id, err := db.AuthCreate(models.DbUserData{
        Email:    userData.Email,
        Username: userData.Username,
        PassHash: passHash,
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

func CheckLoginPass(signInData models.SignInData) (data models.JwtData, err error) {
    var userData models.DbUserData
    passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(signInData.Password)))

    isEmail, _ := regexp.MatchString(`^[\w\-.]+@[\w\-.]+\.[a-z]{2,6}$`, signInData.Login)
    if isEmail {
        userData, err = db.AuthFindByEmailAndPassHash(signInData.Login, passHash)
        if err != nil {
            if err == sql.ErrNoRows {
                return data, fmt.Errorf("incorrect login or password")
            }
            return
        }
    }

    isUsername, _ := regexp.MatchString(`^\w+$`, signInData.Login)
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

func CheckJwt(token string) (models.JwtData, error) {
    data := models.JwtData{}
    err := data.UnMarshal(token, secret)
    return data, err
}
