package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type SignInData struct {
	Login    string `json:"login" example:"test@mail.ru"`
	Password string `json:"password" example:"Qwerty123"`
}

type UserData struct {
	Email    string `json:"email" example:"user_test@test.com"`
	Username string `json:"username" example:"user_test"`
	Password string `json:"password" example:"verysecretpasswordwhichnooneknows"`
}

type DbUserData struct {
	Id       int64    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	PassHash string `json:"pass_hash"`
}

type JwtData struct {
	jwt.StandardClaims
	Id       int64    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (data *JwtData) Marshal(lifetime time.Duration, secret []byte) (string, error) {
	data.StandardClaims.ExpiresAt = time.Now().Add(lifetime).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(secret)
}

func (data *JwtData) UnMarshal(tokenString string, secret []byte) error {
	token, err := jwt.ParseWithClaims(tokenString, data, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if _, ok := token.Claims.(*JwtData); !ok || !token.Valid {
		return err
	}
	return nil
}
