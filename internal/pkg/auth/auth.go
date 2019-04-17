package auth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"2019_1_OPG_plus_2/internal/pkg/db"
	"2019_1_OPG_plus_2/internal/pkg/models"
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

func PasswordHash(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func NewStorage() *Storage {
	return &Storage{}
}

type Storage struct{}

func (*Storage) SignUp(signUpData models.SignUpData) (models.JwtData, error, []string) {
	incorrectFields := signUpData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	id, err := db.AuthCreate(db.AuthData{
		Email:    signUpData.Email,
		Username: signUpData.Username,
		Password: PasswordHash(signUpData.Password),
	})
	if err != nil {
		return models.JwtData{}, err, nil
	}

	return models.JwtData{
		Id:       id,
		Email:    signUpData.Email,
		Username: signUpData.Username,
	}, nil, nil
}

func (*Storage) SignIn(signInData models.SignInData) (data models.JwtData, err error, incorrectFields []string) {
	var userData db.AuthData
	passHash := PasswordHash(signInData.Password)

	isEmail := models.CheckEmail(signInData.Login)
	if isEmail {
		userData, err = db.AuthFindByEmailAndPassHash(signInData.Login, passHash)
		if err != nil {
			if err == models.NotFound {
				return data, models.FieldsError, append(incorrectFields, "password")
			}
			return
		}
	}

	isUsername := !isEmail && models.CheckUsername(signInData.Login)
	if isUsername {
		userData, err = db.AuthFindByUsernameAndPassHash(signInData.Login, passHash)
		if err != nil {
			if err == models.NotFound {
				return data, models.FieldsError, append(incorrectFields, "password")
			}
			return
		}
	}

	if !isEmail && !isUsername {
		return data, models.FieldsError, append(incorrectFields, "login")
	}

	return models.JwtData{
		Id:       userData.Id,
		Email:    userData.Email,
		Username: userData.Username,
	}, nil, nil
}

func (*Storage) UpdateAuth(id int64, userData models.UpdateUserData) (models.JwtData, error, []string) {
	incorrectFields := userData.Check()
	if len(incorrectFields) > 0 {
		return models.JwtData{}, models.FieldsError, incorrectFields
	}

	err := db.AuthUpdateData(db.AuthData{
		Id:       id,
		Email:    userData.Email,
		Username: userData.Username,
	})
	if err != nil {
		return models.JwtData{}, err, nil
	}

	return models.JwtData{
		Id:       id,
		Email:    userData.Email,
		Username: userData.Username,
	}, nil, nil
}

func (*Storage) UpdatePassword(id int64, passwordData models.UpdatePasswordData) (error, []string) {
	incorrectFields := passwordData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	return db.AuthUpdatePassword(id, PasswordHash(passwordData.NewPassword)), nil
}

func (*Storage) RemoveAuth(id int64, removeData models.RemoveUserData) (error, []string) {
	incorrectFields := removeData.Check()
	if len(incorrectFields) > 0 {
		return models.FieldsError, incorrectFields
	}

	err := db.AuthRemove(id, PasswordHash(removeData.Password))
	if err != nil {
		if err.Error() == "incorrect password" {
			return models.FieldsError, append(incorrectFields, "password")
		}
		return err, nil
	}
	return nil, nil
}
