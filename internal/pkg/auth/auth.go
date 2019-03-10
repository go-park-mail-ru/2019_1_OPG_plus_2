package auth

import (
	"crypto/sha256"
	"fmt"
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

	id := create(models.DbUserData{
		Email:    userData.Email,
		Username: userData.Username,
		PassHash: passHash,
	})
	if id == 0 {
		return models.JwtData{}, fmt.Errorf("user already exists")
	}

	return models.JwtData{
		Id:       id,
		Email:    userData.Email,
		Username: userData.Username,
	}, nil
}

func CheckLoginPass(signInData models.SignInData) (models.JwtData, error) {
	var userData *models.DbUserData

	matched, _ := regexp.MatchString(`^[\w\-.]+@[\w\-.]+\.[a-z]{2,6}$`, signInData.Login)
	if matched {
		userData = findByEmail(signInData.Login)
	} else {
		matched, _ = regexp.MatchString(`^\w+$`, signInData.Login)
		if matched {
			userData = findByNickname(signInData.Login)
		} else {
			return models.JwtData{}, fmt.Errorf("incorrect login")
		}
	}

	if userData == nil {
		return models.JwtData{}, fmt.Errorf("incorrect login or password")
	}

	passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(signInData.Password)))
	if passHash != userData.PassHash {
		return models.JwtData{}, fmt.Errorf("incorrect login or password")
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
